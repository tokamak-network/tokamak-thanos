#!/usr/bin/env bash

TOTAL_STEPS=10
STEP=1
SUCCESS="false"

# Detect Operating System
OS_TYPE=$(uname)

# Detect Architecture
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
    ARCH="arm64"
elif [[ "$ARCH" == "armv6l" ]]; then
    ARCH="armv6l"
elif [[ "$ARCH" == "i386" ]]; then
    ARCH="386"
else
    echo "$ARCH is an unsupported architecture."
    exit 1
fi

echo

# MacOS
if [[ "$OS_TYPE" == "Darwin" ]]; then
    OS_TYPE="darwin"
    echo "🚀 Starting package installation for MacOS $ARCH!"

# Linux
elif [[ "$OS_TYPE" == "Linux" ]]; then
    OS_TYPE="linux"
    # Detect the Linux distribution (Ubuntu, Fedora, Arch, etc)
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        OS_NAME=$NAME
        echo "🚀 Starting package installation for Linux/$OS_NAME $ARCH!"
    fi

# Other OS (Windows, BSD, etc.)
else
    echo "$OS_TYPE is an unsupported operating system."
    exit 1
fi

echo

# Check Shell
SHELL_NAME=$(basename "$SHELL")
if [[ "$SHELL_NAME" == "zsh" ]]; then
    echo "The current shell is $SHELL_NAME. The installation will proceed based on $SHELL_NAME."
elif [[ "$SHELL_NAME" == "bash" ]]; then
    echo "The current shell is $SHELL_NAME. The installation will proceed based on $SHELL_NAME."
else
    echo "The current shell is $SHELL_NAME. $SHELL_NAME is an unsupported shell."
    exit 1
fi

# Set Config File
if [ "$SHELL_NAME" = "zsh" ]; then
    CONFIG_FILE="$HOME/.zshrc"
    PROFILE_FILE="$HOME/.zshrc"
elif [ "$SHELL_NAME" = "bash" ]; then
    CONFIG_FILE="$HOME/.bashrc"
    PROFILE_FILE="$HOME/.profile"
fi

echo

# Function to display completion message
function display_completion_message {
    if [[ "$SUCCESS" == "true" ]]; then
        echo ""
        echo "All $TOTAL_STEPS steps are complete."
        echo "Please source your profile to apply changes:"
        echo -e "\033[1;32msource $CONFIG_FILE\033[0m"
        echo ""

        echo "Let's start devnet:"
        echo -e "\033[1;34m\033[1mmake devnet-up\033[0m"
        echo ""
        exit 0
    else
        echo ""
        echo "Installation was interrupted. Completed $((STEP - 1))/$TOTAL_STEPS steps."

        echo ""
        echo "Please source your profile to apply changes:"
        echo -e "\033[1;32msource $CONFIG_FILE\033[0m"
        exit 1
    fi
}

# Use trap to display message on script exit, whether successful or due to an error
trap display_completion_message EXIT
trap "echo 'Process interrupted!'; exit 1" INT

# MacOS specific steps
if [[ "$OS_TYPE" == "darwin" ]]; then

    # 1. Install Homebrew
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Homebrew..."
    if ! command -v brew &> /dev/null; then
        echo "Homebrew not found, installing..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

        if [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
            # Check if the Homebrew configuration is already in the CONFIG_FILE
            if ! grep -Fxq 'export PATH="/usr/local/bin:$PATH"' "$CONFIG_FILE"; then
                # If the configuration is not found, add Homebrew to the current shell session
                {
                    echo ''
                    echo 'export PATH="/usr/local/bin:$PATH"'
                } >> "$CONFIG_FILE"
            fi

            # Check if the Homebrew configuration is already in the PROFILE_FILE
            if ! grep -Fxq 'export PATH="/usr/local/bin:$PATH"' "$PROFILE_FILE"; then
                # If the configuration is not found, add Homebrew to the current shell session
                {
                    echo ''
                    echo 'export PATH="/usr/local/bin:$PATH"'
                } >> "$PROFILE_FILE"
            fi
        fi

        export PATH="/opt/homebrew/bin:$PATH"
    else
        echo "Homebrew is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 2. Install Git
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Git..."
    if ! command -v git &> /dev/null; then
        echo "git not found, installing..."
        brew install git
    else
        echo "git is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 3. Install Xcode Command Line Tools(Inclue make)
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Xcode Command Line Tools..."
    if ! xcode-select -p &> /dev/null; then
        echo "Xcode Command Line Tools not found, installing..."
        xcode-select --install
    else
        echo "Xcode Command Line Tools are already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 4. Install Go (v1.22.6)
    # 4-1. Install Go (v1.22.6)
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Go (v1.22.6)..."
    export PATH="$PATH:/usr/local/go/bin"

    # Save the current Go version
    current_go_version=$(go version 2>/dev/null)

    # Check if the current version is not v1.22.6
    if ! echo "$current_go_version" | grep 'go1.22.6' &>/dev/null ; then

        # If Go is not installed, install Go 1.22.6 directly
        if ! command -v go &> /dev/null; then
            echo "Go not found, installing..."

            if ! command -v curl &> /dev/null; then
                echo "curl not found, installing..."
                brew install curl
            else
                echo "curl is already installed."
            fi

            GO_FILE_NAME="go1.22.6.darwin-${ARCH}.tar.gz"
            GO_DOWNLOAD_URL="https://go.dev/dl/${GO_FILE_NAME}"

            sudo curl -L -o "${GO_FILE_NAME}" "${GO_DOWNLOAD_URL}"

            sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf "${GO_FILE_NAME}"

            # Check if the Go configuration is already in the CONFIG_FILE
            if ! grep -Fxq 'export PATH="$PATH:/usr/local/go/bin"' "$CONFIG_FILE"; then
                # If the configuration is not found, add Go to the current shell session
                {
                    echo ''
                    echo 'export PATH="$PATH:/usr/local/go/bin"'
                } >> "$CONFIG_FILE"
            fi

            # Check if the NVM configuration is already in the PROFILE_FILE
            if ! grep -Fxq 'export PATH=$PATH:/usr/local/go/bin' "$PROFILE_FILE"; then
                # If the configuration is not found, add Go to the current shell session
                {
                    echo ''
                    echo 'export PATH="$PATH:/usr/local/go/bin"'
                } >> "$PROFILE_FILE"
            fi

            export PATH="$PATH:/usr/local/go/bin"

        # If Go is installed and the current version is Go 1.22.6, install GVM.
        else
            # 4-2. Install GVM
            echo "[$STEP/$TOTAL_STEPS] ----- Installing GVM..."
            source ~/.gvm/scripts/gvm
            if ! command -v gvm &> /dev/null; then
                echo "GVM not found, installing..."

                # Install bison for running GVM
                if ! command -v bison &> /dev/null; then
                    echo "bison not found, installing..."
                    apt-get install bison -y
                else
                    echo "bison is already installed."
                fi

                bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)

                # Check if the GVM configuration is already in the CONFIG_FILE
                if ! grep -Fxq 'source ~/.gvm/scripts/gvm' "$CONFIG_FILE"; then

                    # If the configuration is not found, add GVM to the current shell session
                    {
                        echo ''
                        echo 'source ~/.gvm/scripts/gvm'
                    } >> "$CONFIG_FILE"
                fi

                # Check if the GVM configuration is already in the PROFILE_FILE
                if ! grep -Fxq 'source ~/.gvm/scripts/gvm' "$PROFILE_FILE"; then

                    # If the configuration is not found, add GVM to the current shell session
                    {
                        echo ''
                        echo 'source ~/.gvm/scripts/gvm'
                    } >> "$PROFILE_FILE"
                fi

                source ~/.gvm/scripts/gvm
                gvm use system --default
            else
                echo "gvm is already installed."
            fi

            # 4-3. Install Go v1.22.6 using GVM
            echo "[$STEP/$TOTAL_STEPS] ----- Installing Go v1.22.6 using GVM..."
            if ! gvm list | grep 'go1.22.6' &> /dev/null; then
                echo "Go v1.22.6 not found, installing..."
                gvm install go1.22.6
            else
                echo "Go v1.22.6 is already installed."
            fi

            # 4-4. Set Go v1.22.6 as the default version
            echo "[$STEP/$TOTAL_STEPS] ----- Setting Go v1.22.6 as the default version..."
            echo "Switching to Go v1.22.6..."
            gvm use go1.22.6
            echo "Go v1.22.6 is now set as the default version."
        fi
    else
        echo "Go 1.22.6 is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 5. Install Node.js (v20.16.0)
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Node.js (v20.16.0)..."

    # Save the current Node.js version
    current_node_version=$(node -v 2>/dev/null)

    # Check if the current version is not v20.16.0
    if [[ "$current_node_version" != "v20.16.0" ]]; then

        # 5-1. Install NVM
        echo "[$STEP/$TOTAL_STEPS] ----- Installing NVM..."

        # Create NVM directory if it doesn't exist
        export NVM_DIR="$HOME/.nvm"
        mkdir -p "$NVM_DIR"
        HOMEBREW_PREFIX=$(brew --prefix)
        [ -s "$HOMEBREW_PREFIX/opt/nvm/nvm.sh" ] && \. "$HOMEBREW_PREFIX/opt/nvm/nvm.sh"
        [ -s "$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm" ] && \. "$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm"

        if ! command -v nvm &> /dev/null; then
            echo "NVM not found, installing..."
            brew install nvm

            # Check if the NVM configuration is already in the CONFIG_FILE
            if ! grep -Fxq 'export NVM_DIR="$HOME/.nvm"' "$CONFIG_FILE"; then

                # If the configuration is not found, add NVM to the current shell session
                {
                    echo ''
                    echo 'export NVM_DIR="$HOME/.nvm"'
                    echo "[ -s \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\" ] && \. \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\""
                    echo "[ -s \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\" ] && \. \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\""
                } >> "$CONFIG_FILE"
            fi

            # Check if the NVM configuration is already in the PROFILE_FILE
            if ! grep -Fxq 'export NVM_DIR="$HOME/.nvm"' "$PROFILE_FILE"; then

                # If the configuration is not found, add NVM to the current shell session
                {
                    echo ''
                    echo 'export NVM_DIR="$HOME/.nvm"'
                    echo "[ -s \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\" ] && \. \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\""
                    echo "[ -s \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\" ] && \. \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\""
                } >> "$PROFILE_FILE"
            fi

            [ -s "$HOMEBREW_PREFIX/opt/nvm/nvm.sh" ] && \. "$HOMEBREW_PREFIX/opt/nvm/nvm.sh"
            [ -s "$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm" ] && \. "$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm"
        else
            echo "NVM is already installed."
        fi

        # 5-2. Install Node.js v20.16.0 using NVM
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Node.js v20.16.0 using NVM..."
        if ! nvm ls | grep 'v20.16.0' | grep -v 'default' &> /dev/null; then
            echo "Node.js v20.16.0 not found, installing..."
            nvm install v20.16.0
        else
            echo "Node.js v20.16.0 is already installed."
        fi

        # 5-3. Set Node.js v20.16.0 as the default version
        echo "[$STEP/$TOTAL_STEPS] ----- Setting Node.js v20.16.0 as the default version..."
        echo "Switching to Node.js v20.16.0..."
        nvm use v20.16.0
        nvm alias default v20.16.0
        echo "Node.js v20.16.0 is now set as the default version."
    else
        echo "Node.js is already v20.16.0."
    fi

    STEP=$((STEP + 1))
    echo

    # 6. Install Pnpm
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Pnpm..."
    if ! command -v pnpm &> /dev/null; then
        echo "pnpm not found, installing..."
        brew install pnpm
    else
        echo "pnpm is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 7. Install Cargo (v1.78.0)
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Cargo (v1.78.0)..."
    source "$HOME/.cargo/env"
    if ! cargo --version | grep "1.78.0" &> /dev/null; then
        echo "Cargo 1.78.0 not found, installing..."
        curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

        # Check if the Cargo configuration is already in the CONFIG_FILE
        if ! grep -Fq '. "$HOME/.cargo/env"' "$CONFIG_FILE"; then

            # If the configuration is not found, add Cargo to the current shell session
            {
                echo ''
                echo '. "$HOME/.cargo/env"'
            } >> "$CONFIG_FILE"
        fi

        # Check if the Cargo configuration is already in the PROFILE_FILE
        if ! grep -Fq '. "$HOME/.cargo/env"' "$PROFILE_FILE"; then
            # If the configuration is not found, add Cargo to the current shell session
            {
                echo ''
                echo '. "$HOME/.cargo/env"'
            } >> "$PROFILE_FILE"
        fi

        source "$HOME/.cargo/env"
        rustup install 1.78.0
        rustup default 1.78.0
    else
        echo "Cargo 1.78.0 is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 8. Install Docker
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Docker Engine..."
    if ! command -v docker &> /dev/null; then
        echo "Docker not found, installing..."
        brew install --cask docker
    else
        echo "Docker is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 9. Install Foundry using Pnpm
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Foundry using Pnpm..."
    if ! command -v jq &> /dev/null; then
        echo "jq not found, installing..."
        brew install jq
    else
        echo "jq is already installed."
    fi

    if pnpm install:foundry; then
        # Check if the foundry configuration is already in the CONFIG_FILE
        if ! grep -Fq 'export PATH="$PATH:/root/.foundry/bin"' "$CONFIG_FILE"; then

            # If the configuration is not found, add foundry to the current shell session
            {
                echo ''
                echo 'export PATH="$PATH:/root/.foundry/bin"'
            } >> "$CONFIG_FILE"
        fi

        # Check if the foundry configuration is already in the PROFILE_FILE
        if ! grep -Fq 'export PATH="$PATH:/root/.foundry/bin"' "$PROFILE_FILE"; then
            # If the configuration is not found, add foundry to the current shell session
            {
                echo ''
                echo 'export PATH="$PATH:/root/.foundry/bin"'
            } >> "$PROFILE_FILE"
        fi
        export PATH="$PATH:/root/.foundry/bin"
    else
        exit
    fi

    SUCCESS="true"
    echo

# Linux specific steps
elif [[ "$OS_TYPE" == "linux" ]]; then

    # If the operating system is Ubuntu, execute the following commands
    if [[ "$OS_NAME" == "Ubuntu" ]]; then

        if ! command -v sudo &> /dev/null; then
            echo "sudo not found, installing..."
            apt-get install -y sudo
        else
            echo "sudo is already installed."
        fi

        # 1. Update package list
        echo "[$STEP/$TOTAL_STEPS] ----- Updating package list..."
        sudo apt-get update -y

        STEP=$((STEP + 1))
        echo

        # 2. Install Git
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Git..."
        if ! command -v git &> /dev/null; then
            echo "git not found, installing..."
            sudo apt-get install -y git
        else
            echo "git is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 3. Install Build-essential
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Build-essential..."
        if ! dpkg -s build-essential &> /dev/null; then
            echo "Build-essential not found, installing..."
            sudo apt-get install -y build-essential
        else
            echo "Build-essential is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 4. Install Go (v1.22.6)
        # 4-1. Install Go (v1.22.6)
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Go (v1.22.6)..."
        export PATH="$PATH:/usr/local/go/bin"

        # Save the current Go version
        current_go_version=$(go version 2>/dev/null)

        # Check if the current version is not v1.22.6
        if ! echo "$current_go_version" | grep 'go1.22.6' &>/dev/null ; then

            # If Go is not installed, install Go 1.22.6 directly
            if ! command -v go &> /dev/null; then
                echo "Go not found, installing..."

                if ! command -v curl &> /dev/null; then
                    echo "curl not found, installing..."
                    sudo apt-get install -y curl
                else
                    echo "curl is already installed."
                fi

                GO_FILE_NAME="go1.22.6.linux-${ARCH}.tar.gz"
                GO_DOWNLOAD_URL="https://go.dev/dl/${GO_FILE_NAME}"

                sudo curl -L -o "${GO_FILE_NAME}" "${GO_DOWNLOAD_URL}"

                sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf "${GO_FILE_NAME}"

                # Check if the Go configuration is already in the CONFIG_FILE
                if ! grep -Fxq 'export PATH="$PATH:/usr/local/go/bin"' "$CONFIG_FILE"; then
                    # If the configuration is not found, add Go to the current shell session
                    {
                        echo ''
                        echo 'export PATH="$PATH:/usr/local/go/bin"'
                    } >> "$CONFIG_FILE"
                fi

                # Check if the NVM configuration is already in the PROFILE_FILE
                if ! grep -Fxq 'export PATH=$PATH:/usr/local/go/bin' "$PROFILE_FILE"; then
                    # If the configuration is not found, add Go to the current shell session
                    {
                        echo ''
                        echo 'export PATH="$PATH:/usr/local/go/bin"'
                    } >> "$PROFILE_FILE"
                fi

                export PATH="$PATH:/usr/local/go/bin"

            # If Go is installed and the current version is Go 1.22.6, install GVM.
            else
                # 4-2. Install GVM
                echo "[$STEP/$TOTAL_STEPS] ----- Installing GVM..."
                source ~/.gvm/scripts/gvm
                if ! command -v gvm &> /dev/null; then
                    echo "GVM not found, installing..."

                    # Install bison for running GVM
                    if ! command -v bison &> /dev/null; then
                        echo "bison not found, installing..."
                        apt-get install bison -y
                    else
                        echo "bison is already installed."
                    fi

                    bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)

                    # Check if the GVM configuration is already in the CONFIG_FILE
                    if ! grep -Fxq 'source ~/.gvm/scripts/gvm' "$CONFIG_FILE"; then

                        # If the configuration is not found, add GVM to the current shell session
                        {
                            echo ''
                            echo 'source ~/.gvm/scripts/gvm'
                        } >> "$CONFIG_FILE"
                    fi

                    # Check if the GVM configuration is already in the PROFILE_FILE
                    if ! grep -Fxq 'source ~/.gvm/scripts/gvm' "$PROFILE_FILE"; then

                        # If the configuration is not found, add GVM to the current shell session
                        {
                            echo ''
                            echo 'source ~/.gvm/scripts/gvm'
                        } >> "$PROFILE_FILE"
                    fi

                    source ~/.gvm/scripts/gvm
                    gvm use system --default
                else
                    echo "gvm is already installed."
                fi

                # 4-3. Install Go v1.22.6 using GVM
                echo "[$STEP/$TOTAL_STEPS] ----- Installing Go v1.22.6 using GVM..."
                if ! gvm list | grep 'go1.22.6' &> /dev/null; then
                    echo "Go v1.22.6 not found, installing..."
                    gvm install go1.22.6
                else
                    echo "Go v1.22.6 is already installed."
                fi

                # 4-4. Set Go v1.22.6 as the default version
                echo "[$STEP/$TOTAL_STEPS] ----- Setting Go v1.22.6 as the default version..."
                echo "Switching to Go v1.22.6..."
                gvm use go1.22.6
                echo "Go v1.22.6 is now set as the default version."
            fi
        else
            echo "Go 1.22.6 is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 5. Install Node.js (v20.16.0)
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Node.js (v20.16.0)..."

        # Save the current Node.js version
        current_node_version=$(node -v 2>/dev/null)

        # Check if the current version is not v20.16.0
        if [[ "$current_node_version" != "v20.16.0" ]]; then

            # 5-1. Install NVM
            echo "[$STEP/$TOTAL_STEPS] ----- Installing NVM..."

            # Create NVM directory if it doesn't exist
            export NVM_DIR="$HOME/.nvm"
            mkdir -p "$NVM_DIR"
            [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
            [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"

            if ! command -v nvm &> /dev/null; then
                echo "NVM not found, installing..."
                sudo curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash

                # Check if the NVM configuration is already in the CONFIG_FILE
                if ! grep -Fxq 'export NVM_DIR="$HOME/.nvm"' "$CONFIG_FILE"; then

                    # If the configuration is not found, add NVM to the current shell session
                    {
                        echo ''
                        echo 'export NVM_DIR="$HOME/.nvm"'
                        echo "[ -s \"$NVM_DIR/nvm.sh\" ] && \. \"$NVM_DIR/nvm.sh\""
                        echo "[ -s \"$NVM_DIR/bash_completion\" ] && \. \"$NVM_DIR/bash_completion\""
                    } >> "$CONFIG_FILE"
                fi

                # Check if the NVM configuration is already in the PROFILE_FILE
                if ! grep -Fxq 'export NVM_DIR="$HOME/.nvm"' "$PROFILE_FILE"; then

                    # If the configuration is not found, add NVM to the current shell session
                    {
                        echo ''
                        echo 'export NVM_DIR="$HOME/.nvm"'
                        echo "[ -s \"$NVM_DIR/nvm.sh\" ] && \. \"$NVM_DIR/nvm.sh\""
                        echo "[ -s \"$NVM_DIR/bash_completion\" ] && \. \"$NVM_DIR/bash_completion\""
                    } >> "$PROFILE_FILE"
                fi

                [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
                [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"
            else
                echo "NVM is already installed."
            fi

            # 5-2. Install Node.js v20.16.0 using NVM
            echo "[$STEP/$TOTAL_STEPS] ----- Installing Node.js v20.16.0 using NVM..."
            if ! nvm ls | grep 'v20.16.0' | grep -v 'default' &> /dev/null; then
                echo "Node.js v20.16.0 not found, installing..."
                nvm install v20.16.0
            else
                echo "Node.js v20.16.0 is already installed."
            fi

            # 5-3. Set Node.js v20.16.0 as the default version
            echo "[$STEP/$TOTAL_STEPS] ----- Setting Node.js v20.16.0 as the default version..."
            echo "Switching to Node.js v20.16.0..."
            nvm use v20.16.0
            nvm alias default v20.16.0
            echo "Node.js v20.16.0 is now set as the default version."
        else
            echo "Node.js is already v20.16.0."
        fi

        STEP=$((STEP + 1))
        echo

        # 6. Install Pnpm
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Pnpm..."
        export PATH="/root/.local/share/pnpm:$PATH"
        if ! command -v pnpm &> /dev/null; then
            echo "pnpm not found, installing..."
            curl -fsSL https://get.pnpm.io/install.sh | ENV="$CONFIG_FILE" SHELL="$(which "$SHELL_NAME")" "$SHELL_NAME" -

            # Check if the pnpm configuration is already in the CONFIG_FILE
            if ! grep -Fq 'export PATH="/root/.local/share/pnpm:$PATH"' "$CONFIG_FILE"; then

                # If the configuration is not found, add pnpm to the current shell session
                {
                    echo ''
                    echo 'export PATH="/root/.local/share/pnpm:$PATH"'
                } >> "$CONFIG_FILE"
            fi

            # Check if the pnpm configuration is already in the PROFILE_FILE
            if ! grep -Fq 'export PATH="/root/.local/share/pnpm:$PATH"' "$PROFILE_FILE"; then

                # If the configuration is not found, add pnpm to the current shell session
                {
                    echo ''
                    echo 'export PATH="/root/.local/share/pnpm:$PATH"'
                } >> "$PROFILE_FILE"
            fi

            export PATH="/root/.local/share/pnpm:$PATH"
        else
            echo "pnpm is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 7. Install Cargo (v1.78.0)
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Cargo (v1.78.0)..."
        source "$HOME/.cargo/env"
        if ! cargo --version | grep -q "1.78.0" &> /dev/null; then
            echo "Cargo 1.78.0 not found, installing..."
            curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

            # Check if the Cargo configuration is already in the CONFIG_FILE
            if ! grep -Fq '. "$HOME/.cargo/env"' "$CONFIG_FILE"; then

                # If the configuration is not found, add Cargo to the current shell session
                {
                    echo ''
                    echo '. "$HOME/.cargo/env"'
                } >> "$CONFIG_FILE"
            fi

            # Check if the Cargo configuration is already in the PROFILE_FILE
            if ! grep -Fq '. "$HOME/.cargo/env"' "$PROFILE_FILE"; then
                # If the configuration is not found, add Cargo to the current shell session
                {
                    echo ''
                    echo '. "$HOME/.cargo/env"'
                } >> "$PROFILE_FILE"
            fi

            source "$HOME/.cargo/env"
            rustup install 1.78.0
            rustup default 1.78.0
        else
            echo "Cargo 1.78.0 is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 8. Install Docker
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Docker Engine..."
        if ! command -v docker &> /dev/null; then
            echo "Docker not found, installing..."
            sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0
            sudo apt-get install -y gnome-terminal

            # Add Docker's official GPG key:
            sudo apt-get install -y ca-certificates curl
            sudo install -m 0755 -d /etc/apt/keyrings
            sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
            sudo chmod a+r /etc/apt/keyrings/docker.asc

            # Add the repository to Apt sources:
            echo \
              "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
              $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
            # Install the Docker packages.
            sudo apt-get update -y
            sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        else
            echo "Docker is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 9. Install Foundry using Pnpm
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Foundry using Pnpm..."
        if ! command -v jq &> /dev/null; then
            echo "jq not found, installing..."
            sudo apt-get install -y jq
        else
            echo "jq is already installed."
        fi

        if pnpm install:foundry; then
            # Check if the foundry configuration is already in the CONFIG_FILE
            if ! grep -Fq 'export PATH="$PATH:/root/.foundry/bin"' "$CONFIG_FILE"; then

                # If the configuration is not found, add foundry to the current shell session
                {
                    echo ''
                    echo 'export PATH="$PATH:/root/.foundry/bin"'
                } >> "$CONFIG_FILE"
            fi

            # Check if the foundry configuration is already in the PROFILE_FILE
            if ! grep -Fq 'export PATH="$PATH:/root/.foundry/bin"' "$PROFILE_FILE"; then
                # If the configuration is not found, add foundry to the current shell session
                {
                    echo ''
                    echo 'export PATH="$PATH:/root/.foundry/bin"'
                } >> "$PROFILE_FILE"
            fi
            export PATH="$PATH:/root/.foundry/bin"
        else
            exit
        fi

        SUCCESS="true"
        echo

    # If it is an operating system other than Ubuntu, execute the following commands.
    else
        echo "$OS_NAME is an unsupported operating system."
    fi
fi

# Function to check if a command exists and its version if necessary
function check_command_version {
    CMD=$1
    EXPECTED_VERSION=$2
    VERSION_CMD=$3

    if command -v "$CMD" &> /dev/null; then
        CURRENT_VERSION=$($VERSION_CMD 2>&1 | head -n 1)
        if [[ -z "$EXPECTED_VERSION" ]]; then
            if [[ "$CMD" == "forge" || "$CMD" == "cast" || "$CMD" == "anvil" ]]; then
                echo "✅ foundry - $CMD is installed. Current version: $CURRENT_VERSION"
            else
                echo "✅ $CMD is installed. Current version: $CURRENT_VERSION"
            fi
        elif echo "$CURRENT_VERSION" | grep -q "$EXPECTED_VERSION"; then
            echo "✅ $CMD is installed and matches version $EXPECTED_VERSION."
        else
            echo "❌ $CMD is installed but version does not match $EXPECTED_VERSION. Current version: $CURRENT_VERSION"
        fi
    else
        if [[ "$CMD" == "forge" || "$CMD" == "cast" || "$CMD" == "anvil" ]]; then
            echo "❌ foundry - $CMD is not installed."
        else
            echo "❌ $CMD is not installed."
        fi
    fi
}

# Final step: Check installation and versions
echo "Verifying installation and versions..."

# Check Homebrew (for MacOS) - Just check if installed
if [[ "$OS_TYPE" == "darwin" ]]; then
    # Check Homebrew
    check_command_version brew "" "brew --version"

    # Check Git
    check_command_version git "" "git --version"

    # Check Make
    check_command_version make "" "make --version"

    # Check Xcode
    check_command_version xcode-select "" "xcode-select --version"

    # Check Go (Expect version 1.22.6)
    check_command_version go "go1.22.6" "go version"

    # Check Node.js (Expect version 20.16.0)
    check_command_version node "v20.16.0" "node -v"

    # Check Pnpm
    check_command_version pnpm "" "pnpm --version"

    # Check Cargo (Expect version 1.78.0)
    check_command_version cargo "1.78.0" "cargo --version"

    # Check Docker
    check_command_version docker "" "docker --version"

    # Check Foundry
    check_command_version forge "" "forge --version"
    check_command_version cast "" "cast --version"
    check_command_version anvil "" "anvil --version"

    echo "🎉 All required tools are installed and ready to use!"


elif [[ "$OS_TYPE" == "linux" ]]; then
    # Check Git
    check_command_version git "" "git --version"

    # Check Make
    check_command_version make "" "make --version"

    # Check gcc (Instead build-essential)
    check_command_version gcc "" "gcc --version"

    # Check Go (Expect version 1.22.6)
    check_command_version go "go1.22.6" "go version"

    # Check Node.js (Expect version 20.16.0)
    check_command_version node "v20.16.0" "node -v"

    # Check Pnpm
    check_command_version pnpm "" "pnpm --version"

    # Check Cargo (Expect version 1.78.0)
    check_command_version cargo "1.78.0" "cargo --version"

    # Check Docker
    check_command_version docker "" "docker --version"

    # Check Foundry
    check_command_version forge "" "forge --version"
    check_command_version cast "" "cast --version"
    check_command_version anvil "" "anvil --version"

    echo "🎉 All required tools are installed and ready to use!"
fi
