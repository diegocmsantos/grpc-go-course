package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/grpc-go-course/primenumber/primenumberpb"
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

	c := primenumberpb.NewPrimeNumberServiceClient(cc)
	doServerStreaming(c)
}

func doServerStreaming(c primenumberpb.PrimeNumberServiceClient) {
	fmt.Println("Calling the server streaming...")
	stream, err := c.GetPrimeNumbers(context.Background(), &primenumberpb.PrimeNumberRequest{PrimeNumber: 33928293})
	if err != nil {
		log.Fatalf("error calling Prime Number RPC: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of stream")
				break
			}
			fmt.Printf("error receiving message: %v", err)
		}
		fmt.Printf("Response from Prime Number: %v\n", message.GetResult())
	}
}
