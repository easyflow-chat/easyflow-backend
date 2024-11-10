FROM nginx:alpine

# Add the appuser and appgroup
RUN addgroup -g 2000 -S appgroup
RUN adduser -DH -s /sbin/nologin -u 2000 -G appgroup -S appuser

RUN mkdir /app

WORKDIR /app

# Copy the binary from the builder image
COPY --chown=appuser:appgroup ../bin/easyflow-backend ./easyflow-backend
COPY --chown=appuser:appgroup ./nginx.conf /etc/nginx/nginx.conf
COPY --chown=appuser:appgroup ./entrypoint.sh ./entrypoint.sh

# Create the necessary directories with correct permissions
RUN mkdir -p /var/ /var/run /logs/ && \
    chown -R appuser:appgroup /var/ /var/run/ /logs/

# Change the permissions
RUN chmod +x ./easyflow-backend

# Entrypoint script to run both nginx and the Go application
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
