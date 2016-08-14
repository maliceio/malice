#!/usr/bin/env bash
set -e

ZSH_FUNC_DIR="/usr/share/zsh/site-functions"
ZSH_HOMEBREW_FUNC_DIR="/usr/local/share/zsh/site-functions"

if [ -d "$ZSH_HOMEBREW_FUNC_DIR" ]; then
    echo "Installing into ${ZSH_HOMEBREW_FUNC_DIR}..."
    sudo cp ./_malice "$ZSH_HOMEBREW_FUNC_DIR"
    echo "Installed!"
elif [ -d "$ZSH_FUNC_DIR" ]; then
    echo "Installing into ${ZSH_FUNC_DIR}..."
    sudo cp ./_malice "$ZSH_FUNC_DIR"
    echo "Installed! Make sure that ${ZSH_FUNC_DIR} is in your \$fpath."    
else
    echo "Could not find ${ZSH_FUNC_DIR}. Please install manually."
    exit 1
fi
