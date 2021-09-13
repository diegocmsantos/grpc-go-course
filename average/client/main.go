package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grpc-go-course/average/averagepb"
	"google.golang.org/grpc"
)

const port = ":50051"

func main() {
	fmt.Println("Client here")

	cc, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v", err)
		}
	}(cc)

	c := averagepb.NewAverageServiceClient(cc)
	doClientStreaming(c)
}

func doClientStreaming(c averagepb.AverageServiceClient) {
	fmt.Println("Calling the client streaming...")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("error calling Average RPC: %v", err)
	}

	requests := []*averagepb.AverageRequest{
		{Number: 1},
		{Number: 2},
		{Number: 3},
		{Number: 4},
	}

	for _, req := range requests {
		err := stream.Send(req)
		if err != nil {
			fmt.Printf("error while sending request to the server: %v", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("error while receiving response from the server: %v", err)
	}
	fmt.Printf("Response from Average: %v", resp.Average)
}
