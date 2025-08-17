package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/akkahshh24/movieapp/gen"
	metadatatest "github.com/akkahshh24/movieapp/metadata/pkg/testutil"
	movietest "github.com/akkahshh24/movieapp/movie/pkg/testutil"
	ratingtest "github.com/akkahshh24/movieapp/rating/pkg/testutil"

	"github.com/akkahshh24/movieapp/pkg/discovery"
	"google.golang.org/grpc"
)

func startMetadataService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("Starting metadata service on " + metadataServiceAddr)

	// Create a new metadata gRPC server for testing
	handler := metadatatest.NewTestMetadataGRPCServer()
	l, err := net.Listen("tcp", metadataServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Register the metadata service handler
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, handler)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()

	// Register the metadata service with the discovery registry
	id := discovery.GenerateInstanceID(metadataServiceName)
	if err := registry.Register(ctx, id, metadataServiceName, metadataServiceAddr); err != nil {
		panic(err)
	}

	// Periodically report the healthy state of the service
	go func() {
		for {
			if err := registry.ReportHealthyState(id, metadataServiceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return srv
}

func startRatingService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("Starting rating service on " + ratingServiceAddr)

	// Create a new rating gRPC server for testing
	handler := ratingtest.NewTestRatingGRPCServer()
	l, err := net.Listen("tcp", ratingServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Register the rating service handler
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, handler)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()

	// Register the rating service with the discovery registry
	id := discovery.GenerateInstanceID(ratingServiceName)
	if err := registry.Register(ctx, id, ratingServiceName, ratingServiceAddr); err != nil {
		panic(err)
	}

	// Periodically report the healthy state of the service
	go func() {
		for {
			if err := registry.ReportHealthyState(id, ratingServiceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return srv
}

func startMovieService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("Starting movie service on " + movieServiceAddr)

	// Create a new movie gRPC server for testing
	handler := movietest.NewTestMovieGRPCServer(registry)
	l, err := net.Listen("tcp", movieServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Register the movie service handler
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, handler)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()

	// Register the movie service with the discovery registry
	id := discovery.GenerateInstanceID(movieServiceName)
	if err := registry.Register(ctx, id, movieServiceName, movieServiceAddr); err != nil {
		panic(err)
	}

	// Periodically report the healthy state of the service
	go func() {
		for {
			if err := registry.ReportHealthyState(id, movieServiceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return srv
}
