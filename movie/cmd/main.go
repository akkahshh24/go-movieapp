package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/movie/internal/controller/movie"
	metadatagateway "github.com/akkahshh24/movieapp/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/akkahshh24/movieapp/movie/internal/gateway/rating/http"
	grpchandler "github.com/akkahshh24/movieapp/movie/internal/handler/grpc"
	"github.com/akkahshh24/movieapp/movie/pkg/constant"
	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Parse command line flags for the port number.
	// This allows the user to specify which port the movie service should listen on.
	// If no port is specified, it defaults to 8083.
	// This is useful for running multiple instances of the service or for testing purposes.
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)

	// Create a new Consul registry instance.
	// This registry will be used to discover other services in the system.
	// It connects to a Consul agent running on localhost at port 8500.
	// The registry is responsible for service registration and discovery.
	// It allows the movie service to find and communicate with other services like metadata and rating.
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

	// Create an HTTP handler for the movie service.
	// This handler will handle incoming HTTP requests and route them to the appropriate controller methods.
	// h := httphandler.New(ctrl)
	// http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	// if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	// 	panic(err)
	// }

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
