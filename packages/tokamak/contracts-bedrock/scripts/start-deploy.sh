#!/usr/bin/env bash

configFilePath=
environement=
chainID=
deployResultFile=


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
  echo "Start buiding source code"
  source ~/.bashrc
  export PATH=$PATH:/usr/local/go/bin
  cd $projectRoot
  pnpm install
  make submodules
  make cannon-prestate
  make op-node
  cd $projectRoot/packages/tokamak/contracts-bedrock && pnpm build
  cd $projectRoot/packages/tokamak/core-utils && pnpm build
  cd $projectRoot/packages/tokamak/sdk && pnpm build
  cd $currentPWD
}

deployContracts() {
  echo "Start deploying smart contracts"
  echo $DEPLOY_CONFIG_PATH
  export IMPL_SALT=$(openssl rand -hex 32)
  cd $projectRoot/packages/tokamak/contracts-bedrock
  unset DEPLOYMENT_OUTFILE
  if [[ -n "$GAS_PRICE" && "$GAS_PRICE" -gt 0 ]]; then
    forge script scripts/Deploy.s.sol:Deploy -vvv --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive --with-gas-price $GAS_PRICE --resume
  else
    forge script scripts/Deploy.s.sol:Deploy -vvv --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy --non-interactive --resume
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
