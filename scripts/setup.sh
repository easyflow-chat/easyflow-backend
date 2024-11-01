# This is a setup script for the dependencies it needs pipx to be installed

set -e

# Install gloangci-lint
echo "Installling golangci-lint"
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Install precommit
echo "Installing pre-commit"
pipx install pre-commit
pre-commit install
pre-commit run --all


go install github.com/cespare/reflex@latest
echo "Setup completed successfully."
