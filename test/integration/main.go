package main

import (
	"context"
	"log"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/pkg/discovery/memory"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// Service names
	metadataServiceName = "metadata"
	ratingServiceName   = "rating"
	movieServiceName    = "movie"

	// Service addresses
	metadataServiceAddr = "localhost:8081"
	ratingServiceAddr   = "localhost:8082"
	movieServiceAddr    = "localhost:8083"
)

func main() {
	log.Println("Starting the integration test")

	ctx := context.Background()
	registry := memory.NewRegistry()

	// Instantiate our services
	log.Println("Setting up service handlers and clients")

	metadataSrv := startMetadataService(ctx, registry)
	defer metadataSrv.GracefulStop()

	ratingSrv := startRatingService(ctx, registry)
	defer ratingSrv.GracefulStop()

	movieSrv := startMovieService(ctx, registry)
	defer movieSrv.GracefulStop()

	// Setup the test clients for our services
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	metadataConn, err := grpc.NewClient(metadataServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer metadataConn.Close()
	metadataClient := gen.NewMetadataServiceClient(metadataConn)

	ratingConn, err := grpc.NewClient(ratingServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer ratingConn.Close()
	ratingClient := gen.NewRatingServiceClient(ratingConn)

	movieConn, err := grpc.NewClient(movieServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer movieConn.Close()
	movieClient := gen.NewMovieServiceClient(movieConn)

	// Write metadata for an example movie using the metadata service API
	// (the PutMetadata endpoint) and check that the operation does not return any errors.
	log.Println("Metadata service :: PutMetataData :: Writing test metadata")

	m := &gen.Metadata{
		Id:          "the-movie",
		Title:       "The Movie",
		Description: "The Movie, the one and only",
		Director:    "Mr. D",
	}

	if _, err := metadataClient.PutMetadata(ctx, &gen.PutMetadataRequest{Metadata: m}); err != nil {
		log.Fatalf("put metadata: %v", err)
	}

	// Retrieve the metadata for the same movie using the metadata service API
	// (the GetMetadata endpoint) and check it matches the record that we submitted earlier.
	log.Println("Metadata service :: GetMetadata :: Retrieving test metadata")

	getMetadataResp, err := metadataClient.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get metadata: %v", err)
	}

	if diff := cmp.Diff(getMetadataResp.Metadata, m, cmpopts.IgnoreUnexported(gen.Metadata{})); diff != "" {
		log.Fatalf("get metadata after put mismatch: %v", diff)
	}

	// Get the movie details (which should only consist of metadata) for our example movie
	// using the movie service API (the GetMovieDetails endpoint) and make sure the result
	// matches the data that we submitted earlier.
	log.Println("Movie service :: GetMovieDetails :: Getting movie details")

	wantMovieDetails := &gen.MovieDetails{
		Metadata: m,
	}

	getMovieDetailsResp, err := movieClient.GetMovieDetails(ctx, &gen.GetMovieDetailsRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get movie details: %v", err)
	}

	if diff := cmp.Diff(getMovieDetailsResp.MovieDetails, wantMovieDetails, cmpopts.IgnoreUnexported(gen.MovieDetails{}, gen.Metadata{})); diff != "" {
		log.Fatalf("get movie details after put mismatch: %v", err)
	}

	// Write the first rating for our example movie using the rating service API
	// (the PutRating endpoint) and check the operation does not return any errors.
	log.Println("Rating service :: PutRating :: Saving first rating")

	const userID = "user0"
	const recordTypeMovie = "movie"
	firstRating := int32(5)
	if _, err = ratingClient.PutRating(ctx, &gen.PutRatingRequest{
		UserId:      userID,
		RecordId:    m.Id,
		RecordType:  recordTypeMovie,
		RatingValue: firstRating,
	}); err != nil {
		log.Fatalf("put rating: %v", err)
	}

	// Retrieve the initial aggregated rating for our movie using the rating service API
	// (the GetAggregatedRating endpoint) and check that the value matches the one that we just
	// submitted in the previous step.
	log.Println("Rating service :: GetAggregatedRating :: Retrieving initial aggregated rating")

	getAggregatedRatingResp, err := ratingClient.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   m.Id,
		RecordType: recordTypeMovie,
	})
	if err != nil {
		log.Fatalf("get aggreggated rating: %v", err)
	}

	if got, want := getAggregatedRatingResp.RatingValue, float64(5); got != want {
		log.Fatalf("rating mismatch: got %v want %v", got, want)
	}

	// Write the second rating for our example movie using the rating service API
	// and check that the operation does not return any errors.
	log.Println("Rating service :: PutRating :: Saving second rating")

	secondRating := int32(1)
	if _, err = ratingClient.PutRating(ctx, &gen.PutRatingRequest{
		UserId:      userID,
		RecordId:    m.Id,
		RecordType:  recordTypeMovie,
		RatingValue: secondRating,
	}); err != nil {
		log.Fatalf("put rating: %v", err)
	}

	// Retrieve the new aggregated rating for our movie using the rating service API
	// and check that the value reflects the last rating.
	log.Println("Rating service :: GetAggregatedRating :: Retrieving new aggregated rating")

	getAggregatedRatingResp, err = ratingClient.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   m.Id,
		RecordType: recordTypeMovie,
	})
	if err != nil {
		log.Fatalf("get aggreggated rating: %v", err)
	}

	wantRating := float64((firstRating + secondRating) / 2)
	if got, want := getAggregatedRatingResp.RatingValue, wantRating; got != want {
		log.Fatalf("rating mismatch: got %v want %v", got, want)
	}

	// Get the movie details for our example movie and check that the result
	// includes the up-dated rating.
	log.Println("Movie service :: GetMovieDetails :: Getting updated movie details")

	getMovieDetailsResp, err = movieClient.GetMovieDetails(ctx, &gen.GetMovieDetailsRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get movie details: %v", err)
	}

	wantMovieDetails.Rating = wantRating
	if diff := cmp.Diff(getMovieDetailsResp.MovieDetails, wantMovieDetails, cmpopts.IgnoreUnexported(gen.MovieDetails{}, gen.Metadata{})); diff != "" {
		log.Fatalf("get movie details after update mismatch: %v", err)
	}

	log.Println("Integration test execution successful")
}
