package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zkHttp "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	zkReporter reporter.Reporter
	zkTracer   opentracing.Tracer
)

const (
	serviceName     = "zipkin_gin_server"
	serviceEndpoint = "localhost:8080"
	zipkinAddr      = "http://127.0.0.1:9411/api/v2/spans"
)

func initZipkinTracer(engine *gin.Engine) error {
	zkReporter = zkHttp.NewReporter(zipkinAddr)
	endpoint, err := zipkin.NewEndpoint(serviceName, serviceEndpoint)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
		return err
	}
	nativeTracer, err := zipkin.NewTracer(
		zkReporter, zipkin.WithTraceID128Bit(true),
		zipkin.WithLocalEndpoint(endpoint),
	)
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
		return err
	}
	zkTracer = zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(zkTracer)

	// 将tracer注入到gin的中间件中
	engine.Use(func(c *gin.Context) {
		span := zkTracer.StartSpan(c.FullPath())
		defer span.Finish()
		c.Next()
	})
	return nil
}

func main() {
	engine := gin.Default()

	err := initZipkinTracer(engine)
	if err != nil {
		panic(err)
	}
	defer zkReporter.Close()

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	engine.Run(":8080")
}
