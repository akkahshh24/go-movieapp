package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/metadata/internal/controller/metadata"
	grpchandler "github.com/akkahshh24/movieapp/metadata/internal/handler/grpc"
	"github.com/akkahshh24/movieapp/metadata/internal/repository/memory"
	"github.com/akkahshh24/movieapp/metadata/internal/repository/mysql"
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
	log.Printf("Starting the metadata service on port %d", port)

	// Create a new Consul registry
	// This will be used for service discovery.
	registry, err := consul.NewRegistry(cfg.ServiceDiscovery.Consul.Address)
	if err != nil {
		panic(err)
	}

	// Register the metadata service
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

	// Construct DSN in the form: user:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName,
	)

	repo, err := mysql.New(dsn)
	if err != nil {
		panic(err)
	}
	cache := memory.New()
	ctrl := metadata.New(repo, cache)

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
