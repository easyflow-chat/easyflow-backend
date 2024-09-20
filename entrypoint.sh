#!/bin/sh
#Watch for errors and exit immediately
set -e

# Certificates
echo "${CLOUDFLARE_ORIGIN_CERTIFICATE}" > /etc/ssl/backend-easyflow.pem
echo "${CLOUDFLARE_ORIGIN_CA_KEY}" > /etc/ssl/backend-easyflow.key

# Start the Go application in the background
/app/easyflow-backend &

# Start Nginx in the foreground
nginx -g 'daemon off;'