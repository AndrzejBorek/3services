package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	srv "github.com/AndrzejBorek/3services/1st/internal/server"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/generate/json/", srv.LoggingMiddleware(srv.MakeHandler(srv.GenerateJsonHandler)))
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("Service 1 listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
