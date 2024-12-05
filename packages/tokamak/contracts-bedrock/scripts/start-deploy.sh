#!/usr/bin/env bash

configFilePath=
environement=
chainID=
deployResultFile=

OPTSTRING=":c:e:"

projectRoot=`pwd | sed 's%\(.*/tokamak-thanos\)/.*%\1%'`
currentPWD=$(pwd)

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

  export DEPLOY_CONFIG_PATH=$configFilePath
  source $environement

  export chainID=$(jq '.l1ChainID' $configFilePath)
  echo "ChainID: $chainID"

  echo "Project root: $projectRoot"
  echo "Current dir: $currentPWD"
  echo "Config file path: $configFilePath"
  echo "Env file: $environment"

  checkScriptInput
  checkAccount
}

updateSystem() {
  echo "Update system"
  apt update
  DEBIAN_FRONTEND=noninteractive apt install git curl wget cmake jq build-essential -y
}

checkScriptInput() {
  echo "Checking input"
}

checkAccount() {
  echo "Verifying accounts"
}

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
  cd $projectRoot
  pnpm install
  make submodules
  make op-node
  cd $projectRoot/packages/tokamak/contracts-bedrock && pnpm build
}

deployContracts() {
  echo "Start deploying smart contracts"
  echo $DEPLOY_CONFIG_PATH
  cd $projectRoot/packages/tokamak/contracts-bedrock
  forge script scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --legacy
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
  fi

  cd $projectRoot/op-node
  ./bin/op-node genesis l2 \
  --deploy-config $DEPLOY_CONFIG_PATH \
  --l1-deployments $deployResultFile \
  --outfile.l2 $outdir/genesis.json \
  --outfile.rollup $outdir/rollup.json \
  --l1-rpc $L1_RPC_URL
}

updateSystem
handleScriptInput "$@"
installDependency
buildSource
deployContracts
generateL2Genesis
