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
func main() {
	handler := HandlerStruct{}

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start server", err)
	}
}
