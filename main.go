package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HandlerStruct struct{}

func (h HandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type SampleBody struct {
		OK string
	}
	sampleBody := SampleBody{
		OK: "Yes",
	}

	responseBody, err := json.Marshal(sampleBody)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(responseBody)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start the server")
	}
}
