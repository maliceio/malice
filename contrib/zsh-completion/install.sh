#!/bin/bash
set -e

ZSH_FUNC_DIR="~/.yadr/zsh/prezto/modules/completion/external/src"

if [ -d "$ZSH_FUNC_DIR" ]; then
    echo "Installing into ${ZSH_FUNC_DIR}..."
    sudo cp ./_malice "$ZSH_FUNC_DIR"
    echo "Installed! Make sure that ${ZSH_FUNC_DIR} is in your \$fpath."
else
    echo "Could not find ${ZSH_FUNC_DIR}. Please install manually."
    exit 1
fi
