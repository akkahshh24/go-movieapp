package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/discovery/consul"
	"github.com/akkahshh24/movieapp/rating/internal/controller/rating"
	httphandler "github.com/akkahshh24/movieapp/rating/internal/handler/http"
	"github.com/akkahshh24/movieapp/rating/internal/repository/memory"
	"github.com/akkahshh24/movieapp/rating/pkg/constant"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting the rating service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	// Register the rating service
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(constant.ServiceNameRating)
	if err := registry.Register(ctx, instanceID, constant.ServiceNameRating, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	// Periodically report the healthy state of the service
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, constant.ServiceNameRating); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Deregister the service on exit
	defer registry.Deregister(ctx, instanceID, constant.ServiceNameRating)

	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
