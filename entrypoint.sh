#!/bin/sh

# Start the Go application in the background
/app/easyflow-backend &

# Start Nginx in the foreground
nginx -g 'daemon off;'