package main

import (
	"log"
	"net/http"

	"github.com/akkahshh24/movieapp/metadata/internal/controller/metadata"
	httphandler "github.com/akkahshh24/movieapp/metadata/internal/handler/http"
	"github.com/akkahshh24/movieapp/metadata/internal/repository/memory"
)

func main() {
	log.Println("Starting the metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
