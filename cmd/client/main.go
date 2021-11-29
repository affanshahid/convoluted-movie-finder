package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/affanshahid/convoluted-movie-finder/core"
	"github.com/affanshahid/convoluted-movie-finder/rpc/pb"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	url          = flag.String("host", "localhost:50051", "URL of the movie-finder server")
	genreId      = flag.Int64("g", 28, "ID of the genre to search into")
	startDateStr = flag.String("s", "2021-11-12", "Starting date of the search interval")
	endDateStr   = flag.String("e", "2021-11-13", "Ending date of the search interval")
	revenue      = flag.Int64("r", 1000, "Revenue threshold")
	operator     = flag.Int("o", int(core.OpGt), "Operator to use when comparing revenue, 0: <, 1: ==, 2: >")
)

func init() {
	flag.Parse()
}

func main() {
	conn, err := grpc.Dial(*url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMovieClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startTime, err := time.Parse("2006-01-02", *startDateStr)
	if err != nil {
		panic(err)
	}

	endTime, err := time.Parse("2006-01-02", *endDateStr)
	if err != nil {
		panic(err)
	}

	if *operator > 2 || *operator < 0 {
		panic("Operator must be 0,1 or 2")
	}

	r, err := client.FetchGenrePeriodDetails(ctx, &pb.GenrePeriodDetailsRequest{
		GenreId:              *genreId,
		StartDate:            timestamppb.New(startTime),
		EndDate:              timestamppb.New(endTime),
		Revenue:              *revenue,
		RevenueCheckOperator: pb.GenrePeriodDetailsRequest_Operator(*operator),
	})
	if err != nil {
		panic(err)
	}

	spew.Dump(r)
}
