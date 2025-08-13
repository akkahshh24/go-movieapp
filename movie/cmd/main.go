package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/movie/internal/controller/movie"
	metadatagateway "github.com/akkahshh24/movieapp/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/akkahshh24/movieapp/movie/internal/gateway/rating/http"
	grpchandler "github.com/akkahshh24/movieapp/movie/internal/handler/grpc"
	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/discovery/consul"
	"github.com/akkahshh24/movieapp/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

func main() {
	f, err := os.Open("default.yaml")
	if err != nil {
		panic(err)
	}

	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	port := cfg.API.Port
	log.Printf("Starting the movie service on port %d", port)

	// Create a new Consul registry instance.
	// This registry will be used to discover other services in the system.
	// It connects to a Consul agent running on localhost at port 8500.
	// The registry is responsible for service registration and discovery.
	// It allows the movie service to find and communicate with other services like metadata and rating.
	registry, err := consul.NewRegistry(cfg.ServiceDiscovery.Consul.Address)
	if err != nil {
		panic(err)
	}

	// Register the movie service
	ctx := context.Background()
	serviceName := model.ServiceName(cfg.ServiceDiscovery.Name)
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%s:%d", serviceName, port)); err != nil {
		panic(err)
	}

	// Periodically report the healthy state of the service
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Deregister the service on exit
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)

	// Create a gRPC server and register the movie service.
	// This server will listen for incoming gRPC requests on the specified port.
	// It will use the movie controller to handle requests related to movie operations.
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
