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
    wget https://go.dev/dl/go1.21.13.linux-amd64.tar.gz
    tar xvzf go1.21.13.linux-amd64.tar.gz
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
  retryCommand "pnpm install" "Installing dependencies" || return 1

  # Initialize submodules
  retryCommand "make submodules" "Initializing submodules" || return 1

  # Build cannon prestate
  retryCommand "make cannon-prestate" "Building cannon prestate" || return 1

  # Build op-node
  retryCommand "make op-node" "Building op-node" || return 1

  # Build contracts-bedrock
  echo "Building contracts-bedrock..."
  cd $projectRoot/packages/tokamak/contracts-bedrock

  # Ensure forge is available
  if ! command -v forge &> /dev/null; then
    echo "Error: forge command not found. Please install Foundry first."
    return 1
  fi

  # Clean and build contracts with retry logic
  retryCommand "forge clean && forge build" "Building contracts" || return 1

  # Verify forge artifacts exist and wait if necessary
  verifyForgeArtifacts "forge-artifacts" || return 1

  # Wait for file system sync before building TypeScript packages
  waitForFileSystem

  # Build TypeScript packages in dependency order
  echo "Building core-utils..."
  cd $projectRoot/packages/tokamak/core-utils

  # Additional wait to ensure modules are properly synced
  waitForFileSystem

  # Build core-utils with retry logic
  retryCommand "pnpm build" "Building core-utils" || return 1

  # Build SDK
  echo "Building SDK..."
  cd $projectRoot/packages/tokamak/sdk

  # Additional wait to ensure modules are properly synced
  waitForFileSystem

  # Build SDK with retry logic
  retryCommand "pnpm build" "Building SDK" || return 1

  cd $currentPWD
  echo "All source code built successfully!"
}

deployContracts() {
  echo "Start deploying smart contracts"
  echo $DEPLOY_CONFIG_PATH
  export IMPL_SALT=$(openssl rand -hex 32)
  cd $projectRoot/packages/tokamak/contracts-bedrock
  unset DEPLOYMENT_OUTFILE
  if [[ -n "$GAS_PRICE" && "$GAS_PRICE" -gt 0 ]]; then
    forge script scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive --with-gas-price $GAS_PRICE
  else
    forge script scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive
  fi
  cd $currentPWD
}

resumeDeployContracts() {
  echo "Resume deploying smart contracts"
  echo $DEPLOY_CONFIG_PATH
  cd $projectRoot/packages/tokamak/contracts-bedrock

  if [[ -n "$GAS_PRICE" && "$GAS_PRICE" -gt 0 ]]; then
    forge script scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive --with-gas-price $GAS_PRICE --resume
  else
    forge script scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive --resume
  fi
  cd $currentPWD
}

generateL2Genesis() {
  echo "Generate L2 genesis"
  deployResultFile=$projectRoot/packages/tokamak/contracts-bedrock/deployments/$(printf "%d-deploy.json" "$chainID")
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
  $projectRoot/op-node/bin/op-node genesis l2 \
  --deploy-config $DEPLOY_CONFIG_PATH \
  --l1-deployments $deployResultFile \
  --outfile.l2 $outdir/genesis.json \
  --outfile.rollup $outdir/rollup.json \
  --l1-rpc $L1_RPC_URL

  echo "Genesis file: $outdir/genesis.json"
  echo "Rollup file: $outdir/rollup.json"
  cd $currentPWD
}

main() {
  case $1 in
    install)
      echo "Install softwares..."
      updateSystem
      installDependency
      buildSource
      ;;
    build)
      echo "Build..."
      buildSource
      ;;
    deploy)
      echo "Deploying smart contracts..."
      shift
      handleScriptInput "$@"
      deployContracts
      ;;
    redeploy)
      echo "Redeploying smart contracts..."
      shift
      handleScriptInput "$@"
      resumeDeployContracts
      ;;
    generate)
      echo "Generate rollup and genesis config for L2..."
      shift
      handleScriptInput "$@"
      generateL2Genesis
      ;;
    all)
      echo "Setup from scratch"
      updateSystem
      shift
      handleScriptInput "$@"
      installDependency
      buildSource
      deployContracts
      generateL2Genesis
      ;;
    *)
      echo "Usage: $0 {install|build|deploy|generate|all}"
      exit 1
      ;;
esac
}

main "$@"
exit 0
