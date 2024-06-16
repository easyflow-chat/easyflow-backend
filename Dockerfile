# Build stage
FROM golang:1.22.4 as builder

WORKDIR /app
COPY . .

#Clean up the go modules
RUN go mod tidy

# Lint the Go application
RUN gofmt -w .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/easyflow-backend ./src

# Production stage
FROM nginx:alpine as production

# Add the appuser and appgroup
RUN addgroup -g 2000 -S appgroup
RUN adduser -DH -s /sbin/nologin -u 2000 -G appgroup -S appuser

RUN mkdir /app
RUN chown -R appuser:appgroup /app

WORKDIR /app

# Copy the binary from the builder image
COPY --chown=appuser:appgroup --from=builder /app/bin/easyflow-backend ./easyflow-backend
COPY --chown=appuser:appgroup --from=builder /app/nginx.conf /etc/nginx/nginx.conf
COPY --chown=appuser:appgroup --from=builder /app/entrypoint.sh ./entrypoint.sh

# Create the necessary directories with correct permissions
RUN mkdir -p /var/cache/nginx/client_temp && \
    chown -R appuser:appgroup /var/cache/nginx && \
    mkdir -p /var/run/ && \
    chown -R appuser:appgroup /var/run/


# Change the user to appuser
USER appuser

# Metadata
LABEL org.opencontainers.image.authors="nico.benninger43@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/easyflow-chat/easyflow-backend"
LABEL org.opencontainers.image.title="Easyflow Backend"
LABEL org.opencontainers.image.description="Backend for Easyflow chat application"

# Change the permissions
RUN chmod +x ./easyflow-backend

# Entrypoint script to run both nginx and the Go application
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
