package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/openzipkin/zipkin-go"
	zkhttpmw "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/reporter"
	zkhttpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	zkReporter reporter.Reporter
)

const (
	serviceName     = "zipkin_http_server"
	serviceEndpoint = "localhost:8080"
	zipkinAddr      = "http://127.0.0.1:9411/api/v2/spans"
)

func Pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pong")
	return
}

func initMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", Pong)

	zkReporter = zkhttpreporter.NewReporter(zipkinAddr)
	endpoint, err := zipkin.NewEndpoint(serviceName, serviceEndpoint)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	tracer, err := zipkin.NewTracer(
		zkReporter, zipkin.WithTraceID128Bit(true),
		zipkin.WithLocalEndpoint(endpoint),
	)

	zkMiddleware := zkhttpmw.NewServerMiddleware(tracer)
	return zkMiddleware(mux)
}

func main() {
	mux := initMux()

	http.ListenAndServe(":8080", mux)
}
