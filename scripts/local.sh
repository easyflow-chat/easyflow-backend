# setup a local environment for development
# What we need:
# - Local mariadb instance hosted on 3306 accessible by localhost:3306
# - Local Adminer instance hosted on 8080 accessible by localhost:8080
if [ -z "$1" ]; then
    echo "Usage: $0 <start|stop>"
    exit 1
fi

start_env() {
    echo "Starting local environment..."
    echo "Ensuring there is not already a mariadb container running..."
    if [ "$(docker ps -q -f name=mariadb)" ]; then
        echo "Stopping existing mariadb container..."
        docker stop mariadb
    fi

    echo "Starting mariadb container..."
    docker run --name mariadb -e MYSQL_ROOT_PASSWORD=devel -p 3306:3306 -d mariadb:latest
    if [ $? -ne 0 ]; then
        echo "Failed to start mariadb container."
        exit 1
    fi

    echo "Ensuring there is not already an adminer container running..."
    if [ "$(docker ps -q -f name=adminer)" ]; then
        echo "Stopping existing adminer container..."
        docker stop adminer
    fi

    echo "Starting adminer container..."
    docker run --name adminer -p 8080:8080 --link mariadb:db -d adminer:latest
    if [ $? -ne 0 ]; then
        echo "Failed to start adminer container."
        exit 1
    fi

    echo "Local environment started successfully."
}

stop_env() {
    echo "Stopping local environment..."
    echo "Stopping mariadb container..."
    docker stop mariadb
    echo "Stopping adminer container..."
    docker stop adminer
    echo "Local environment stopped successfully."
}

case "$1" in
    start)
        start_env
        ;;
    stop)
        stop_env
        ;;
    *)
        echo "Usage: $0 <start|stop>"
        exit 1
        ;;
esac

