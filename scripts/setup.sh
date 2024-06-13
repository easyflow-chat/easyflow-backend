#!/bin/bash

# Check if pre-commit is installed and install it if it's not
if ! command -v pre-commit &> /dev/null
then
    echo "Installing pre-commit..."
    # check if brew is installed install it over brew
    if command -v brew &> /dev/null
    then
        brew install pre-commit
    else
        pip install pre-commit
    fi
fi

# Install the pre-commit hooks
echo "Installing pre-commit hooks..."
pre-commit install
go install github.com/cespare/reflex@latest
echo "Setup completed successfully."
