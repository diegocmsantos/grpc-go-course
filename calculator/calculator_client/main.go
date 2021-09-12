package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc-go-course/calculator/calculatorpb"
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

	c := calculatorpb.NewCalculatorServiceClient(cc)
	resp, err := c.Sum(context.Background(), &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{A: 18, B: 5},
	})
	if err != nil {
		log.Fatalf("error calling Greet RPC: %v", err)
	}
	fmt.Printf("Response from Greet: %v", resp.Sum)
}