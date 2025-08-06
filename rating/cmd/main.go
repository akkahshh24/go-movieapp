package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/discovery/consul"
	"github.com/akkahshh24/movieapp/rating/internal/controller/rating"
	grpchandler "github.com/akkahshh24/movieapp/rating/internal/handler/grpc"
	"github.com/akkahshh24/movieapp/rating/internal/ingester/kafka"
	"github.com/akkahshh24/movieapp/rating/internal/repository/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "rating"

func main() {
	// Take the port from command line arguments
	// Default to 8082 if not provided.
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting the rating service on port %d", port)

	// Create a new Consul registry
	// This will be used for service discovery.
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	// Register the rating service.
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	// Periodically report the healthy state of the service.
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	// Deregister the service on exit.
	defer registry.Deregister(ctx, instanceID, serviceName)

	// Create the repository and controller.
	repo := memory.New()
	ingester, err := kafka.NewIngester("localhost", "rating", "ratings")
	if err != nil {
		log.Fatalf("failed to initialize ingester: %v", err)
	}
	ctrl := rating.New(repo, ingester)
	// Start the consumer to ingest rating events.
	// This will listen to the Kafka topic and process incoming rating events.
	if err := ctrl.StartIngestion(ctx); err != nil {
		log.Fatalf("failed to start ingestion: %v", err)
	}

	// Create the gRPC handler and register it with the gRPC server.
	// This handler will implement the gRPC service methods.
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterRatingServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
