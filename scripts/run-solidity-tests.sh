#!/bin/bash

set -e

# kill any existing `titand` processes
process_name="titand"
if pgrep -x "$process_name" >/dev/null; then
    # Get the PID of the process
    pid=$(pgrep -x "$process_name")

    # Kill the process
    echo "Killing process $process_name with PID $pid"
    kill $pid
fi

# Install nvm if not present
export NVM_DIR="$HOME/.nvm"
if [ ! -d "$NVM_DIR" ]; then
    echo "Installing NVM..."
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

    # Source NVM immediately after installation
    export NVM_DIR="$HOME/.nvm"
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
else
    # Source NVM if it exists
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
fi

# Set Node.js version to 18
echo "Setting Node.js version to 18..."
nvm use 18 || nvm install 18
npm install -g yarn

node -v

export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin

make install

cd tests/e2e/solidity

if command -v yarn &>/dev/null; then
    yarn install
else
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
    sudo apt update && sudo apt install yarn
    yarn install
fi

# shellcheck disable=SC2068
yarn test --network titan $@

# kill any existing `titand` processes
process_name="titand"
if pgrep -x "$process_name" >/dev/null; then
    # Get the PID of the process
    pid=$(pgrep -x "$process_name")

    # Kill the process
    echo "Killing process $process_name with PID $pid"
    kill $pid
fi
