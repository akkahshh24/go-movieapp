package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akkahshh24/movieapp/metadata/internal/controller/metadata"
	httphandler "github.com/akkahshh24/movieapp/metadata/internal/handler/http"
	"github.com/akkahshh24/movieapp/metadata/internal/repository/memory"
	"github.com/akkahshh24/movieapp/metadata/pkg/constant"
	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/discovery/consul"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the metadata service on port %d", port)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	// Register the metadata service
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(constant.ServiceNameMetadata)
	if err := registry.Register(ctx, instanceID, constant.ServiceNameMetadata, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	// Periodically report the healthy state of the service
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, constant.ServiceNameMetadata); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Deregister the service on exit
	defer registry.Deregister(ctx, instanceID, constant.ServiceNameMetadata)

	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
