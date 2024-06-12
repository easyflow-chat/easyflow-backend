#!/bin/bash

# Check if pre-commit is installed and install it if it's not
if ! command -v pre-commit &> /dev/null
then
    echo "Installing pre-commit..."
    pip install pre-commit
fi

# Install the pre-commit hooks
echo "Installing pre-commit hooks..."
pre-commit install

echo "Setup completed successfully."
