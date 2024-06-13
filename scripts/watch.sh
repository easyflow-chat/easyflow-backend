#!/bin/bash

# Start watching for file changes and run the development script
reflex -r '^src/.*\.go$' -s -- sh scripts/dev.sh &

# Capture the PID of the reflex process
#REFLEX_PID=$!

# Pipe the Docker container logs to the terminal
#docker logs -f docker-easyflow-backend-1 &

# Wait for the Reflex process to complete
#wait $REFLEX_PID
