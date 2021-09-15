package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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

	//doUnary(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	resp, err := c.Sum(context.Background(), &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{A: 18, B: 5},
	})
	if err != nil {
		log.Fatalf("error calling Greet RPC: %v", err)
	}
	fmt.Printf("Response from Greet: %v", resp.Sum)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Square Root unary RPC...")

	resp, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: -30})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Printf("Expected error with code [%s] and message [%s]", respErr.Code(), respErr.Message())
			return
		}
		log.Fatalf("error receiving response from the server: %v\n", err)
		return
	}
	fmt.Printf("Response from Square Root RPC: %.2f\n", resp.GetNumberRoot())
}