#!/bin/bash

# Build the Go application
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/easyflow-backend ./src

# Check if the build was successful
if [ $? -ne 0 ]; then
  echo "Build failed, stopping execution."
  exit 1
fi

# Stop and remove old containers
docker compose -f ./docker/docker-compose.yml down

# If the old containers are successfully removed, proceed with Docker Compose
docker compose -f ./docker/docker-compose.yml up --build -d

docker system prune -f