TOTAL_STEPS=10

echo "🚀 Starting package installation for MacOS!"

echo

# Check Shell
if [ -n "$ZSH_VERSION" ]; then
    SHELL_NAME="zsh"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_NAME="bash"
fi

# Set Config File
if [ "$SHELL_NAME" = "zsh" ]; then
    CONFIG_FILE="$HOME/.zshrc"
elif [ "$SHELL_NAME" = "bash" ]; then
    CONFIG_FILE="$HOME/.profile"
fi

echo "The current shell is $SHELL_NAME. The installation will proceed based on $SHELL_NAME."

echo

# 1. Install Homebrew
echo "[1/$TOTAL_STEPS] ----- Installing Homebrew..."
if ! command -v brew &> /dev/null; then
    echo "Homebrew not found, installing..."
    if [ "$SHELL_NAME" = "zsh" ]; then
        /bin/zsh -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    elif [ "$SHELL_NAME" = "bash" ]; then
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    fi
else
    echo "Homebrew is already installed."
fi

echo

# 2. Install Git
echo "[2/$TOTAL_STEPS] ----- Installing Git..."
if ! command -v git &> /dev/null; then
    echo "git not found, installing..."
    brew install git
else
    echo "git is already installed."
fi

echo

# 3. Install Make
echo "[3/$TOTAL_STEPS] ----- Installing Make..."
if ! command -v make &> /dev/null; then
    echo "make not found, installing..."
    brew install make
else
    echo "make is already installed."
fi

echo

# 4. Install Build-essential
echo "[4/$TOTAL_STEPS] ----- Installing Xcode Command Line Tools..."
if ! xcode-select -p &> /dev/null; then
    echo "Xcode Command Line Tools not found, installing..."
    xcode-select --install
else
    echo "Xcode Command Line Tools are already installed."
fi

echo

# 5. Install Go (v1.22.6)
echo "[5/$TOTAL_STEPS] ----- Installing Go (v1.22.6)..."
if ! go version | grep "go1.22.6" &> /dev/null; then
    echo "Go 1.22.6 not found, installing..."
    brew install go@1.22
    brew link --force --overwrite go@1.22
else
    echo "Go 1.22.6 is already installed."
fi

echo

# 6. Install Node.js (v20.16.0)
echo "[6/$TOTAL_STEPS] ----- Installing Node.js (v20.16.0)..."

## 6-1. Install NVM
if ! command -v nvm &> /dev/null; then
    echo "NVM not found, installing..."
    brew install nvm

    # Create NVM directory if it doesn't exist
    export NVM_DIR="$HOME/.nvm"
    mkdir -p "$NVM_DIR"

    # Extract the HOMEBREW_PREFIX path
    HOMEBREW_PREFIX=$(brew config | grep 'HOMEBREW_PREFIX' | awk '{print $2}')

    # Check if the NVM configuration is already in the .zshrc file
    if ! grep -Fxq 'export NVM_DIR="$HOME/.nvm"' "$CONFIG_FILE"; then

        # If the configuration is not found, add NVM to the current shell session
        {
            echo ''
            echo 'export NVM_DIR="$HOME/.nvm"'
            echo "[ -s \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\" ] && \. \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\""
            echo "[ -s \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\" ] && \. \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\""
        } >> "$CONFIG_FILE"

        if [ "$SHELL_NAME" = "zsh" ]; then
            source ~/.zshrc
        elif [ "$SHELL_NAME" = "bash" ]; then
            source ~/.profile
        fi
    fi

else
    echo "NVM is already installed."
fi

## 6-2. Install Node.js v20.16.0 using NVM
echo "[6/$TOTAL_STEPS] ----- Installing Node.js v20.16.0 using NVM..."
if ! nvm ls | grep "v20.16.0" &> /dev/null; then
    echo "Node.js v20.16.0 not found, installing..."
    nvm install v20.16.0
else
    echo "Node.js v20.16.0 is already installed."
fi

## 6-3. Set Node.js v20.16.0 as the default version
echo "[6/$TOTAL_STEPS] ----- Setting Node.js v20.16.0 as the default version..."

# Save the current Node.js version
current_version=$(node -v 2>/dev/null)

# Check if the current version is not v20.16.0
if ! echo "$current_version" | grep "v20.16.0" &> /dev/null; then
    echo "Current Node.js version: $current_version"
    echo "Switching to Node.js v20.16.0..."
    nvm use v20.16.0
    nvm alias default v20.16.0
    echo "Node.js v20.16.0 is now set as the default version."
else
    echo "Node.js is already v20.16.0."
fi

echo

# 7. Install Pnpm
echo "[7/$TOTAL_STEPS] ----- Installing Pnpm..."
if ! brew list pnpm &> /dev/null; then
    echo "pnpm not found, installing..."
    brew install pnpm
else
    echo "pnpm is already installed."
fi

echo

# 8. Install Cargo (v1.78.0)
echo "[8/$TOTAL_STEPS] ----- Installing Cargo (v1.78.0)..."
if ! cargo --version | grep "1.78.0" &> /dev/null; then
    echo "Cargo 1.78.0 not found, installing..."
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    source $HOME/.cargo/env
    rustup install 1.78.0
    rustup default 1.78.0
else
    echo "Cargo 1.78.0 is already installed."
fi

echo

# 9. Install Docker Engine
echo "[9/$TOTAL_STEPS] ----- Installing Docker Engine..."
if ! command -v docker &> /dev/null; then
    echo "Docker not found, installing..."
    brew install --cask docker
    # Add Docker alias to shell config file if not already present
    if ! grep -Fxq 'alias docker="/Applications/Docker.app/Contents/Resources/bin/docker"' "$CONFIG_FILE"; then
        echo 'alias docker="/Applications/Docker.app/Contents/Resources/bin/docker"' >> "$CONFIG_FILE"

        if [ "$SHELL_NAME" = "zsh" ]; then
            source ~/.zshrc
        elif [ "$SHELL_NAME" = "bash" ]; then
            source ~/.profile
        fi
    fi

else
    echo "Docker is already installed."
fi

echo

# 10. Install Foundry using Pnpm
echo "[10/$TOTAL_STEPS] ----- Installing Foundry using Pnpm..."
if ! brew list libusb &> /dev/null; then
    echo "libusb not found, installing..."
    brew install libusb
else
    echo "libusb is already installed."
fi

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

echo

echo "All $TOTAL_STEPS steps are complete."
