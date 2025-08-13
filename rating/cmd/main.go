package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/discovery/consul"
	"github.com/akkahshh24/movieapp/pkg/model"
	"github.com/akkahshh24/movieapp/rating/internal/cache/memory"
	"github.com/akkahshh24/movieapp/rating/internal/controller/rating"
	grpchandler "github.com/akkahshh24/movieapp/rating/internal/handler/grpc"
	"github.com/akkahshh24/movieapp/rating/internal/ingester/kafka"
	"github.com/akkahshh24/movieapp/rating/internal/repository/mysql"
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
	log.Printf("Starting the rating service on port %d", port)

	// Create a new Consul registry
	// This will be used for service discovery.
	registry, err := consul.NewRegistry(cfg.ServiceDiscovery.Consul.Address)
	if err != nil {
		panic(err)
	}

	// Register the rating service.
	ctx := context.Background()
	serviceName := model.ServiceName(cfg.ServiceDiscovery.Name)
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%s:%d", serviceName, port)); err != nil {
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

	// Create and in-memory or mysql repository.
	// Here we are using MySQL as the repository.
	// You can switch to an in-memory repository for testing purposes.
	// Construct DSN in the form: user:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName,
	)

	repo, err := mysql.New(dsn)
	if err != nil {
		panic(err)
	}
	cache := memory.New()

	ingester, err := kafka.NewIngester(cfg.MessageQueue.Address, cfg.MessageQueue.GroupID, cfg.MessageQueue.Topic)
	if err != nil {
		log.Fatalf("failed to initialize ingester: %v", err)
	}

	ctrl := rating.New(repo, cache, ingester)

	// Start the consumer to ingest rating events.
	// This will listen to the Kafka topic and process incoming rating events.
	go func() {
		if err := ctrl.StartIngestion(ctx); err != nil {
			log.Fatalf("failed to start ingestion: %v", err)
		}
	}()

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
