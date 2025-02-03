package main

import (
	"log"
	"net/http"
	"reviewllama/internal/handler"
)

func main() {
	http.HandleFunc("/review", handler.ReviewPost)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
