#!/usr/bin/env bash

configFilePath=
environement=
chainID=
deployResultFile=

# Build configuration
MAX_RETRIES=3
RETRY_DELAY=5
BUILD_TIMEOUT=300

OPTSTRING=":c:e:"

projectRoot=`pwd | sed 's%\(.*/tokamak-thanos\)/.*%\1%'`
currentPWD=$(pwd)
configDir=$projectRoot/packages/tokamak/contracts-bedrock/deploy-config/tmp

handleScriptInput() {
  echo "Check script input"
  while getopts ${OPTSTRING} opt; do
    case ${opt} in
      c)
        configFilePath=$OPTARG
        ;;
      e)
        echo "Has env file"
        environement=$OPTARG
        ;;
      :)
        echo "Error: Option -$OPTARG requires an argument." >&2
        exit 1
        ;;
      *)
        echo "Error: Invalid option -$OPTARG" >&2
        exit 1
        ;;
    esac
  done

  if [ ! -d "$configDir" ]; then
    mkdir -p "$configDir"
    echo " Temp config directory '$configDir' created."
  fi

  cp $configFilePath $configDir/config.json
  export DEPLOY_CONFIG_PATH=${configDir}/config.json

  source $environement

  export chainID=$(jq '.l1ChainID' $configFilePath)
  echo "ChainID: $chainID"

  echo "Project root: $projectRoot"
  echo "Current dir: $currentPWD"
  echo "Config file path: $DEPLOY_CONFIG_PATH"
  echo "Env file: $environment"

  # checkScriptInput
  # checkAccount
}

updateSystem() {
  echo "Update system"
  apt update
  DEBIAN_FRONTEND=noninteractive apt install git curl wget cmake jq build-essential -y
}

# Helper functions for robust building
waitForFileSystem() {
  echo "Waiting for file system to sync..."
  sleep 2
  sync
}

verifyForgeArtifacts() {
  local artifacts_dir="$1"
  local max_wait=30
  local wait_count=0

  echo "Verifying forge artifacts in $artifacts_dir..."

  while [ $wait_count -lt $max_wait ]; do
    if [ -d "$artifacts_dir" ] && [ -n "$(ls -A "$artifacts_dir" 2>/dev/null)" ]; then
      echo "Forge artifacts verified successfully"
      return 0
    fi

    echo "Waiting for forge artifacts... ($((wait_count + 1))/$max_wait)"
    sleep 1
    wait_count=$((wait_count + 1))
  done

  echo "Error: Forge artifacts not found after $max_wait seconds"
  return 1
}

retryCommand() {
  local command="$1"
  local description="$2"
  local max_retries=${3:-$MAX_RETRIES}
  local retry_count=0

  echo "Executing: $description"

  while [ $retry_count -lt $max_retries ]; do
    echo "Attempt $((retry_count + 1)) of $max_retries..."

    if eval "$command"; then
      echo "$description completed successfully!"
      return 0
    else
      echo "$description failed. Retrying..."
      retry_count=$((retry_count + 1))
      sleep $RETRY_DELAY
    fi
  done

  echo "Error: $description failed after $max_retries attempts"
  return 1
}

# checkScriptInput() {
#   echo "Checking input"
# }

# checkAccount() {
#   echo "Verifying accounts"
# }

installDependency() {
  installGo
  installJSPackages
  installFoundry
}

installGo() {
  echo "Install Go"
  export PATH=$PATH:/usr/local/go/bin
  if ! go version &> /dev/null; then
    echo "Installing Go..."
    cd /root
    wget https://go.dev/dl/go1.23.8.linux-amd64.tar.gz
    tar xvzf go1.23.8.linux-amd64.tar.gz
    mv go /usr/local
    export PATH=$PATH:/usr/local/go/bin
    echo PATH=$PATH:/usr/local/go/bin >> ~/.bashrc
    echo "   Installed $(go version)"
    cd $currentPWD
  else
    echo "   Found $(go version)"
  fi
}

installJSPackages() {
  curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
  DEBIAN_FRONTEND=noninteractive apt install -y nodejs
  DEBIAN_FRONTEND=noninteractive apt install -y npm
  npm install -g pnpm
  curl -sL https://dl.yarnpkg.com/debian/pubkey.gpg | gpg --dearmor | tee /usr/share/keyrings/yarnkey.gpg >/dev/null
  echo "deb [signed-by=/usr/share/keyrings/yarnkey.gpg] https://dl.yarnpkg.com/debian stable main" | tee /etc/apt/sources.list.d/yarn.list
  apt update && apt install yarn
}

installFoundry() {
  echo "Install Foundry"
  if ! forge --version &> /dev/null; then
    curl -L https://foundry.paradigm.xyz | bash
    mv /root/.foundry/bin/* /usr/bin/
    source ~/.bashrc
    foundryup
    echo "export PATH=/root/.foundry/bin:$PATH" >> /root/.bashrc
    source /root/.bashrc
    export PATH=/root/.foundry/bin:$PATH
    echo $PATH
    echo "   Installed $(forge --version)"
  else
    echo "   Found $(forge --version)"
  fi
}

buildSource() {
  echo "Start building source code"
  source ~/.bashrc
  export PATH=$PATH:/usr/local/go/bin
  cd $projectRoot

  # Install dependencies
  echo "Installing dependencies..."
  if ! retryCommand "pnpm install" "Installing dependencies"; then
    echo "❌ Error: Failed to install dependencies after $MAX_RETRIES attempts"
    return 1
  fi

  # Initialize submodules
  echo "Initializing submodules..."
  if ! retryCommand "make submodules" "Initializing submodules"; then
    echo "❌ Error: Failed to initialize submodules after $MAX_RETRIES attempts"
    return 1
  fi

  # Build cannon prestate
  echo "Building cannon prestate..."
  if ! retryCommand "make cannon-prestate" "Building cannon prestate"; then
    echo "❌ Error: Failed to build cannon prestate after $MAX_RETRIES attempts"
    return 1
  fi

  # Build op-node
  echo "Building op-node..."
  if ! retryCommand "make op-node" "Building op-node"; then
    echo "❌ Error: Failed to build op-node after $MAX_RETRIES attempts"
    return 1
  fi

  # Build contracts-bedrock
  echo "Building contracts-bedrock..."
  cd $projectRoot/packages/tokamak/contracts-bedrock

  # Ensure forge is available
  if ! command -v forge &> /dev/null; then
    echo "❌ Error: forge command not found. Please install Foundry first."
    return 1
  fi

  # Clean and build contracts with retry logic
  echo "Cleaning and building Solidity contracts..."
  if ! retryCommand "forge clean && forge build" "Building contracts"; then
    echo "❌ Error: Failed to build contracts after $MAX_RETRIES attempts"
    return 1
  fi

  # Verify forge artifacts exist and wait if necessary
  if ! verifyForgeArtifacts "forge-artifacts"; then
    echo "❌ Error: Forge artifacts not found"
    return 1
  fi

  # Wait for file system sync before building TypeScript packages
  waitForFileSystem

  # Build TypeScript packages in dependency order
  echo "Building core-utils..."
  cd $projectRoot/packages/tokamak/core-utils

  # Additional wait to ensure modules are properly synced
  waitForFileSystem

  # Build core-utils with retry logic
  if ! retryCommand "pnpm build" "Building core-utils"; then
    echo "❌ Error: Failed to build core-utils after $MAX_RETRIES attempts"
    return 1
  fi

  # Verify core-utils build output
  if [ ! -f "dist/index.js" ]; then
    echo "❌ Error: core-utils build output not found at dist/index.js"
    echo "Listing dist directory:"
    ls -la dist/ 2>/dev/null || echo "dist directory does not exist"
    return 1
  fi
  echo "✅ core-utils build verified: dist/index.js exists"

  # Build SDK
  echo "Building SDK..."
  cd $projectRoot/packages/tokamak/sdk

  # Additional wait to ensure modules are properly synced
  waitForFileSystem

  # Reinstall SDK dependencies to ensure workspace symlinks are correct
  echo "Ensuring SDK workspace dependencies are properly linked..."
  if ! retryCommand "pnpm install --prefer-offline" "Reinstalling SDK dependencies"; then
    echo "❌ Error: Failed to reinstall SDK dependencies after $MAX_RETRIES attempts"
    return 1
  fi

  # Build SDK with retry logic
  if ! retryCommand "pnpm build" "Building SDK"; then
    echo "❌ Error: Failed to build SDK after $MAX_RETRIES attempts"
    return 1
  fi

  # Verify SDK build output and workspace symlinks
  if [ ! -f "dist/index.js" ]; then
    echo "❌ Error: SDK build output not found at dist/index.js"
    echo "Listing dist directory:"
    ls -la dist/ 2>/dev/null || echo "dist directory does not exist"
    return 1
  fi

  # Verify core-utils symlink in SDK node_modules
  if [ ! -e "node_modules/@tokamak-network/core-utils/dist/index.js" ]; then
    echo "⚠️  Warning: core-utils symlink not properly set in SDK node_modules"
    echo "Checking symlink target:"
    ls -la node_modules/@tokamak-network/core-utils/ 2>/dev/null || echo "core-utils not found in node_modules"

    # Try to fix by reinstalling
    echo "Attempting to fix workspace symlinks..."
    cd $projectRoot
    if ! retryCommand "pnpm install --force" "Force reinstalling all dependencies"; then
      echo "❌ Error: Failed to fix workspace symlinks after $MAX_RETRIES attempts"
      return 1
    fi

    # Verify again
    cd $projectRoot/packages/tokamak/sdk
    if [ ! -e "node_modules/@tokamak-network/core-utils/dist/index.js" ]; then
      echo "❌ Error: Failed to restore core-utils symlink after force reinstall"
      return 1
    fi
  fi
  echo "✅ SDK build verified: dist/index.js exists and workspace symlinks are correct"

  cd $currentPWD
  echo "✅ All source code built successfully!"
  return 0
}

deployContracts() {
  echo "Start deploying smart contracts"
  echo $DEPLOY_CONFIG_PATH
  export IMPL_SALT=$(openssl rand -hex 32)
  cd $projectRoot/packages/tokamak/contracts-bedrock
  unset DEPLOYMENT_OUTFILE

  # Pass FPS parameters to op-challenger
  echo "Set FPS params to op-challenger..."
  if [ -n "$CHALLENGE_WINDOW" ]; then
    echo "export CHALLENGE_WINDOW=$CHALLENGE_WINDOW" >> .env
  fi
  if [ -n "$CHALLENGER_PRIVATE_KEY" ]; then
    echo "export CHALLENGER_PRIVATE_KEY=$CHALLENGER_PRIVATE_KEY" >> .env
  fi

  echo "Deploying contracts to L1..."
  local deploy_result
  if [[ -n "$GAS_PRICE" && "$GAS_PRICE" -gt 0 ]]; then
    forge script scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive --with-gas-price $GAS_PRICE
    deploy_result=$?
  else
    forge script scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive
    deploy_result=$?
  fi

  cd $currentPWD

  if [ $deploy_result -ne 0 ]; then
    echo "❌ Error: Contract deployment failed with exit code $deploy_result"
    return 1
  fi

  echo "✅ Contract deployment completed successfully"
  return 0
}

resumeDeployContracts() {
  echo "Resume deploying smart contracts"
  echo $DEPLOY_CONFIG_PATH
  cd $projectRoot/packages/tokamak/contracts-bedrock

  echo "Resuming contract deployment..."
  local deploy_result
  if [[ -n "$GAS_PRICE" && "$GAS_PRICE" -gt 0 ]]; then
    forge script scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive --with-gas-price $GAS_PRICE --resume
    deploy_result=$?
  else
    forge script scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive --resume
    deploy_result=$?
  fi

  cd $currentPWD

  if [ $deploy_result -ne 0 ]; then
    echo "❌ Error: Resume contract deployment failed with exit code $deploy_result"
    return 1
  fi

  echo "✅ Resume contract deployment completed successfully"
  return 0
}

generateL2Genesis() {
  echo "Generate L2 genesis"
  deployResultFile=$projectRoot/packages/tokamak/contracts-bedrock/deployments/$(printf "%d-deploy.json" "$chainID")

  # Check if deployment file exists
  if [ ! -f "$deployResultFile" ]; then
    echo "❌ Error: Deployment file not found: $deployResultFile"
    return 1
  fi

  echo "Deployment file found: $deployResultFile"
  cat $deployResultFile

  export outdir=$projectRoot/build

  if [ ! -d "$outdir" ]; then
    mkdir -p "$outdir"
    echo "Directory '$outdir' created."
  else
    echo "Directory '$outdir' already exists."
    rm -rf $outdir/*
  fi

  cd $projectRoot
  echo "Generating L2 genesis and rollup configuration..."

  if ! $projectRoot/op-node/bin/op-node genesis l2 \
    --deploy-config $DEPLOY_CONFIG_PATH \
    --l1-deployments $deployResultFile \
    --outfile.l2 $outdir/genesis.json \
    --outfile.rollup $outdir/rollup.json \
    --l1-rpc $L1_RPC_URL; then
    echo "❌ Error: Failed to generate L2 genesis and rollup configuration"
    cd $currentPWD
    return 1
  fi

  # Verify generated files
  if [ ! -f "$outdir/genesis.json" ]; then
    echo "❌ Error: Genesis file was not created: $outdir/genesis.json"
    cd $currentPWD
    return 1
  fi

  if [ ! -f "$outdir/rollup.json" ]; then
    echo "❌ Error: Rollup file was not created: $outdir/rollup.json"
    cd $currentPWD
    return 1
  fi

  echo "✅ Genesis file: $outdir/genesis.json"
  echo "✅ Rollup file: $outdir/rollup.json"
  cd $currentPWD
  return 0
}

main() {
  case $1 in
    install)
      echo "Install softwares..."
      updateSystem
      installDependency || { echo "❌ Installation failed"; exit 1; }
      buildSource || { echo "❌ Build failed"; exit 1; }
      echo "✅ Installation completed successfully"
      ;;
    build)
      echo "Build..."
      buildSource || { echo "❌ Build failed"; exit 1; }
      echo "✅ Build completed successfully"
      ;;
    deploy)
      echo "Deploying smart contracts..."
      shift
      handleScriptInput "$@" || { echo "❌ Script input handling failed"; exit 1; }
      deployContracts || { echo "❌ Contract deployment failed"; exit 1; }
      echo "✅ Contract deployment completed successfully"
      ;;
    redeploy)
      echo "Redeploying smart contracts..."
      shift
      handleScriptInput "$@" || { echo "❌ Script input handling failed"; exit 1; }
      resumeDeployContracts || { echo "❌ Resume deployment failed"; exit 1; }
      echo "✅ Resume deployment completed successfully"
      ;;
    generate)
      echo "Generate rollup and genesis config for L2..."
      shift
      handleScriptInput "$@" || { echo "❌ Script input handling failed"; exit 1; }
      generateL2Genesis || { echo "❌ Genesis generation failed"; exit 1; }
      echo "✅ Genesis generation completed successfully"
      ;;
    all)
      echo "Setup from scratch"
      updateSystem
      shift
      handleScriptInput "$@" || { echo "❌ Script input handling failed"; exit 1; }
      installDependency || { echo "❌ Installation failed"; exit 1; }
      buildSource || { echo "❌ Build failed"; exit 1; }
      deployContracts || { echo "❌ Contract deployment failed"; exit 1; }
      generateL2Genesis || { echo "❌ Genesis generation failed"; exit 1; }
      echo "✅ All steps completed successfully"
      ;;
    *)
      echo "Usage: $0 {install|build|deploy|generate|all}"
      exit 1
      ;;
  esac
}

main "$@"
exit_code=$?

if [ $exit_code -eq 0 ]; then
  echo "✅ Script completed successfully"
else
  echo "❌ Script failed with exit code $exit_code"
fi

exit $exit_code
