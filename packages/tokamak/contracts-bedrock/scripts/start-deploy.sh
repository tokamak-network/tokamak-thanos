#!/usr/bin/env bash

configFilePath=
environement=

projectRoot=`pwd | sed 's%\(.*/tokamak-thanos\)/.*%\1%'`
currentPWD=$(pwd)

handleScriptInput() {
  echo "Check script input"
  while getopts "c:e" opt;
  do
    case ${opt} in
      c)
        configFilePath=$OPTARG
        ;;
      e)
        blockNum=$OPTARG
        ;;
    esac
  done

  checkScriptInput
  checkAccount

  echo "Project root: $projectRoot"
  echo "Current dir: $currentPWD"
  echo "Config file path: $configFilePath"
  echo "Env file: $environment"
}

updateSystem() {
  source $environement
  echo "Update system"
  apt update
  DEBIAN_FRONTEND=noninteractive apt install git curl wget cmake jq -y
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
  if ! go version &> /dev/null; then
    echo "Installing Go..."
    cd /root
    wget https://go.dev/dl/go1.21.linux-amd64.tar.gz
    tar xvzf go1.21.linux-amd64.tar.gz
    cp go/bin/go /usr/bin/go
    mv go /usr/local
    echo export GOROOT=/usr/local/go >> ~/.bashrc
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
  cd $projectRoot
  pnpm install
  make build
}

updateSystem
handleScriptInput
installDependency
buildSource
