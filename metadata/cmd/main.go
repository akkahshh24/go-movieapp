package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/metadata/internal/controller/metadata"
	grpchandler "github.com/akkahshh24/movieapp/metadata/internal/handler/grpc"
	"github.com/akkahshh24/movieapp/metadata/internal/repository/memory"
	"github.com/akkahshh24/movieapp/metadata/internal/repository/mysql"
	"github.com/akkahshh24/movieapp/metadata/pkg/constant"
	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Take the port from command line arguments
	// Default to 8081 if not provided.
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the metadata service on port %d", port)

	// Create a new Consul registry
	// This will be used for service discovery.
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

	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}
	cache := memory.New()
	ctrl := metadata.New(repo, cache)

	/* HTTP handler setup
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
	*/

	// gRPC handler setup
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
