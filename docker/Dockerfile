FROM alpine:latest

# Install Dockerize
RUN apk add --no-cache openssl
ENV DOCKERIZE_VERSION v0.7.0
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

WORKDIR /app
COPY bin/easyflow-backend ./easyflow-backend
COPY example.env ./.env
RUN chmod +x easyflow-backend


CMD ["dockerize", "-wait", "tcp://mariadb:3306", "-timeout", "30s", "./easyflow-backend"]
