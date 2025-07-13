package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akkahshh24/movieapp/movie/internal/controller/movie"
	metadatagateway "github.com/akkahshh24/movieapp/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/akkahshh24/movieapp/movie/internal/gateway/rating/http"
	httphandler "github.com/akkahshh24/movieapp/movie/internal/handler/http"
	"github.com/akkahshh24/movieapp/movie/pkg/constant"
	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/discovery/consul"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	// Register the movie service
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(constant.ServiceNameMovie)
	if err := registry.Register(ctx, instanceID, constant.ServiceNameMovie, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	// Periodically report the healthy state of the service
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, constant.ServiceNameMovie); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Deregister the service on exit
	defer registry.Deregister(ctx, instanceID, constant.ServiceNameMovie)

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
