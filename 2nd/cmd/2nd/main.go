package main

import (
	"log"
	"net/http"

	srv "github.com/AndrzejBorek/3services/2nd/internal/server"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/1st/", srv.LoggingMiddleware(srv.MakeHandler(srv.FirstEndpointHandler)))
	mux.Handle("/2nd/", srv.LoggingMiddleware(srv.MakeHandler(srv.SecondEndpointHandler)))
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("Service 2 listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))

}
