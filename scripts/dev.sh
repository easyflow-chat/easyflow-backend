CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/easyflow-backend ./src
docker compose -f ./Docker/docker-compose.yml up --build -d
