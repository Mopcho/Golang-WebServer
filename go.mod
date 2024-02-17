module github.com/Mopcho/Golang-WebServer

go 1.21.6

require github.com/go-chi/chi/v5 v5.0.11

require golang.org/x/crypto v0.19.0 // indirect

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/joho/godotenv v1.5.1
	internal/database v1.0.0
)

replace internal/database => ./internal/database
