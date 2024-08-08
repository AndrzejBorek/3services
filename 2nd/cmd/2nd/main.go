package main

import (
	"log"
	"net/http"

	srv "github.com/AndrzejBorek/3services/2nd/internal/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/1stendpoint", srv.LoggingMiddleware(srv.MakeHandler(srv.FirstEndpointHandler)))
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))

}
