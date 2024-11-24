module easyflow-websocket

go 1.23.3

replace github.com/easyflow-chat/easyflow-backend/lib/database v0.0.0 => ./../lib/database

replace github.com/easyflow-chat/easyflow-backend/lib/logger v0.0.0 => ./../lib/logger

replace github.com/easyflow-chat/easyflow-backend/lib/jwt v0.0.0 => ./../lib/jwt

require (
	github.com/easyflow-chat/easyflow-backend/lib/database v0.0.0
	github.com/easyflow-chat/easyflow-backend/lib/jwt v0.0.0
	github.com/easyflow-chat/easyflow-backend/lib/logger v0.0.0
	github.com/gorilla/websocket v1.5.3
	github.com/joho/godotenv v1.5.1
	gorm.io/gorm v1.25.12
)

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.14.0 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
)
