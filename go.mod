module github.com/Mopcho/Golang-WebServer

go 1.21.6

require github.com/go-chi/chi/v5 v5.0.11

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.19.0 // indirect
)

require internal/database v1.0.0

replace internal/database => ./internal/database
