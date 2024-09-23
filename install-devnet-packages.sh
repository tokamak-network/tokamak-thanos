#!/usr/bin/env bash

TOTAL_STEPS=10
STEP=1

# Detect Operating System
OS_TYPE=$(uname)

# Detect Architecture
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "aarch64" ]]; then
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
    if [[ $STEP -gt $TOTAL_STEPS ]]; then
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

    # 3. Install Make
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Make..."
    if ! command -v make &> /dev/null; then
        echo "make not found, installing..."
        brew install make
    else
        echo "make is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 4. Install Xcode Command Line Tools
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Xcode Command Line Tools..."
    if ! xcode-select -p &> /dev/null; then
        echo "Xcode Command Line Tools not found, installing..."
        xcode-select --install
    else
        echo "Xcode Command Line Tools are already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 5. Install Go (v1.22.6)
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Go (v1.22.6)..."
    if ! go version | grep "go1.22.6" &> /dev/null; then
        echo "Go 1.22.6 not found, installing..."
        brew install go@1.22
        brew link --force --overwrite go@1.22
    else
        echo "Go 1.22.6 is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 6. Install Node.js (v20.16.0)

    # 6-1. Install NVM
    echo "[$STEP/$TOTAL_STEPS] ----- Installing NVM..."

    # Create NVM directory if it doesn't exist
    export NVM_DIR="$HOME/.nvm"
    mkdir -p "$NVM_DIR"

    # Extract the HOMEBREW_PREFIX path
    HOMEBREW_PREFIX=$(brew --prefix)
    [ -s "$HOMEBREW_PREFIX/opt/nvm/nvm.sh" ] && \. "$HOMEBREW_PREFIX/opt/nvm/nvm.sh"
    [ -s "$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm" ] && \. "$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm"
    hash -r

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

        if [ "$SHELL_NAME" = "zsh" ]; then
            source ~/.zshrc
        elif [ "$SHELL_NAME" = "bash" ]; then
            source ~/.bashrc
            source ~/.profile
        fi
    else
        echo "NVM is already installed."
    fi

    # 6-2. Install Node.js v20.16.0 using NVM
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Node.js v20.16.0 using NVM..."
    if ! nvm ls | grep "v20.16.0" &> /dev/null; then
        echo "Node.js v20.16.0 not found, installing..."
        nvm install v20.16.0
    else
        echo "Node.js v20.16.0 is already installed."
    fi

    # 6-3. Set Node.js v20.16.0 as the default version
    echo "[$STEP/$TOTAL_STEPS] ----- Setting Node.js v20.16.0 as the default version..."

    # Save the current Node.js version
    current_version=$(node -v 2>/dev/null)

    # Check if the current version is not v20.16.0
    if [[ "$current_version" != "v20.16.0" ]]; then
        echo "Current Node.js version: $current_version"
        echo "Switching to Node.js v20.16.0..."
        nvm use v20.16.0
        nvm alias default v20.16.0
        echo "Node.js v20.16.0 is now set as the default version."

        if [ "$SHELL_NAME" = "zsh" ]; then
            source ~/.zshrc
        elif [ "$SHELL_NAME" = "bash" ]; then
            source ~/.bashrc
            source ~/.profile
        fi
    else
        echo "Node.js is already v20.16.0."
    fi

    STEP=$((STEP + 1))
    echo

    # 7. Install Pnpm
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Pnpm..."
    if ! command -v pnpm &> /dev/null; then
        echo "pnpm not found, installing..."
        brew install pnpm

        if [ "$SHELL_NAME" = "zsh" ]; then
            source ~/.zshrc
        elif [ "$SHELL_NAME" = "bash" ]; then
            source ~/.bashrc
            source ~/.profile
        fi
    else
        echo "pnpm is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 8. Install Cargo (v1.78.0)
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Cargo (v1.78.0)..."
    source "$HOME/.cargo/env"
    if ! cargo --version | grep "1.78.0" &> /dev/null; then
        echo "Cargo 1.78.0 not found, installing..."
        curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
        rustup install 1.78.0
        rustup default 1.78.0

        if [ "$SHELL_NAME" = "zsh" ]; then
            source ~/.zshrc
        elif [ "$SHELL_NAME" = "bash" ]; then
            source ~/.bashrc
            source ~/.profile
        fi
    else
        echo "Cargo 1.78.0 is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 9. Install Docker
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Docker Engine..."
    if ! command -v docker &> /dev/null; then
        echo "Docker not found, installing..."
        brew install --cask docker

        if [ "$SHELL_NAME" = "zsh" ]; then
            source ~/.zshrc
        elif [ "$SHELL_NAME" = "bash" ]; then
            source ~/.bashrc
            source ~/.profile
        fi
    else
        echo "Docker is already installed."
    fi

    STEP=$((STEP + 1))
    echo

    # 10. Install Foundry using Pnpm
    echo "[$STEP/$TOTAL_STEPS] ----- Installing Foundry using Pnpm..."
    if ! command -v jq &> /dev/null; then
        echo "jq not found, installing..."
        brew install jq
    else
        echo "jq is already installed."
    fi

    if command -v pnpm &> /dev/null; then
        pnpm install:foundry
    else
        echo "Pnpm is not installed. Skipping Foundry installation."
    fi

    STEP=$((STEP + 1))
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

        # 3. Install Make
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Make..."
        if ! command -v make &> /dev/null; then
            echo "make not found, installing..."
            sudo apt-get install -y make
        else
            echo "make is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 4. Install Build-essential
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Build-essential..."
        if ! dpkg -s build-essential &> /dev/null; then
            echo "Build-essential not found, installing..."
            sudo apt-get install -y build-essential
        else
            echo "Build-essential is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 5. Install Go (v1.22.6)
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Go (v1.22.6)..."
        export PATH="$PATH:/usr/local/go/bin"
        hash -r
        if ! go version | grep "go1.22.6" &> /dev/null; then
            echo "Go 1.22.6 not found, installing..."

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

            if [ "$SHELL_NAME" = "zsh" ]; then
                source ~/.zshrc
            elif [ "$SHELL_NAME" = "bash" ]; then
                source ~/.bashrc
                source ~/.profile
            fi
        else
            echo "Go 1.22.6 is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 6. Install Node.js (v20.16.0)

        # 6-1. Install NVM
        echo "[$STEP/$TOTAL_STEPS] ----- Installing NVM..."

        # Create NVM directory if it doesn't exist
        export NVM_DIR="$HOME/.nvm"
        mkdir -p "$NVM_DIR"
        [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
        [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"
        hash -r

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

            if [ "$SHELL_NAME" = "zsh" ]; then
                source ~/.zshrc
            elif [ "$SHELL_NAME" = "bash" ]; then
                source ~/.bashrc
                source ~/.profile
            fi
        else
            echo "NVM is already installed."
        fi

        # 6-2. Install Node.js v20.16.0 using NVM
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Node.js v20.16.0 using NVM..."
        if ! nvm ls | grep "v20.16.0" &> /dev/null; then
            echo "Node.js v20.16.0 not found, installing..."
            nvm install v20.16.0
        else
            echo "Node.js v20.16.0 is already installed."
        fi

        # 6-3. Set Node.js v20.16.0 as the default version
        echo "[$STEP/$TOTAL_STEPS] ----- Setting Node.js v20.16.0 as the default version..."
        # Save the current Node.js version
        current_version=$(node -v 2>/dev/null)

        # Check if the current version is not v20.16.0
        if [[ "$current_version" != "v20.16.0" ]]; then
            echo "Current Node.js version: $current_version"
            echo "Switching to Node.js v20.16.0..."
            nvm use v20.16.0
            nvm alias default v20.16.0
            echo "Node.js v20.16.0 is now set as the default version."

            if [ "$SHELL_NAME" = "zsh" ]; then
                source ~/.zshrc
            elif [ "$SHELL_NAME" = "bash" ]; then
                source ~/.bashrc
                source ~/.profile
            fi
        else
            echo "Node.js is already v20.16.0."
        fi

        STEP=$((STEP + 1))
        echo

        # 7. Install Pnpm
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Pnpm..."
        export PNPM_HOME="$HOME/.local/share/pnpm"
        export PATH="$PNPM_HOME:$PATH"
        hash -r
        if ! command -v pnpm &> /dev/null; then
            echo "pnpm not found, installing..."
            curl -fsSL https://get.pnpm.io/install.sh | ENV="$CONFIG_FILE" SHELL="$(which "$SHELL_NAME")" "$SHELL_NAME" -

            if [ "$SHELL_NAME" = "zsh" ]; then
                source ~/.zshrc
            elif [ "$SHELL_NAME" = "bash" ]; then
                source ~/.bashrc
                source ~/.profile
            fi
        else
            echo "pnpm is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 8. Install Cargo (v1.78.0)
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Cargo (v1.78.0)..."
        source "$HOME/.cargo/env"
        if ! cargo --version | grep -q "1.78.0" &> /dev/null; then
            echo "Cargo 1.78.0 not found, installing..."
            curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
            rustup install 1.78.0
            rustup default 1.78.0
            source "$HOME/.cargo/env"
            if [ "$SHELL_NAME" = "zsh" ]; then
                source ~/.zshrc
            elif [ "$SHELL_NAME" = "bash" ]; then
                source ~/.bashrc
                source ~/.profile
            fi
        else
            echo "Cargo 1.78.0 is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 9. Install Docker
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Docker Engine..."
        if ! command -v docker &> /dev/null; then
            echo "Docker not found, installing..."
            sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0
            sudo apt-get install -y gnome-terminal

            # Add Docker's official GPG key:
            sudo apt-get update -y
            sudo apt-get install -y ca-certificates curl
            sudo install -m 0755 -d /etc/apt/keyrings
            sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
            sudo chmod a+r /etc/apt/keyrings/docker.asc
            # Add the repository to Apt sources:

            echo \
              "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
              $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
            # Install the Docker packages.
            sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

            if [ "$SHELL_NAME" = "zsh" ]; then
                source ~/.zshrc
            elif [ "$SHELL_NAME" = "bash" ]; then
                source ~/.bashrc
                source ~/.profile
            fi
        else
            echo "Docker is already installed."
        fi

        STEP=$((STEP + 1))
        echo

        # 10. Install Foundry using Pnpm
        echo "[$STEP/$TOTAL_STEPS] ----- Installing Foundry using Pnpm..."
        if ! command -v jq &> /dev/null; then
            echo "jq not found, installing..."
            sudo apt-get install -y jq
        else
            echo "jq is already installed."
        fi

        if ! pnpm -y install:foundry; then
            exit
        fi

        STEP=$((STEP + 1))
        echo

    # If it is an operating system other than Ubuntu, execute the following commands.
    else
        echo "$OS_NAME is an unsupported operating system."
    fi
fi
