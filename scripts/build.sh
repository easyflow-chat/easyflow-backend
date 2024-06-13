# get os
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
# if linux or darwin
if [ $OS = "linux" ] || [ $OS = "darwin" ]; then
  # Build the Go application
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/easyflow-backend ./src
  # Check if the build was successful
  if [ $? -ne 0 ]; then
    echo "Build failed, stopping execution."
    exit 1
  fi
fi
# if windows
if [ $OS = "windows" ]; then
  # Build the Go application
  CGO_ENABLED=0 GOOS=windows go build -a -installsuffix cgo -o ./bin/easyflow-backend.exe ./src
  # Check if the build was successful
  if [ $? -ne 0 ]; then
    echo "Build failed, stopping execution."
    exit 1
  fi
fi

#CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/easyflow-backend ./src