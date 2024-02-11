module github.com/Mopcho/Golang-WebServer

go 1.21.6

require github.com/go-chi/chi/v5 v5.0.11

require github.com/google/uuid v1.6.0 // indirect

require internal/database v1.0.0

replace internal/database => ./internal/database
