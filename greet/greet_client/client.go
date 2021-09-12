package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/grpc-go-course/greet/greetpb"
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

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)
	doServerStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	resp, err := c.Greet(context.Background(), &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Diego",
			LastName: "Maia",
		},
	})
	if err != nil {
		log.Fatalf("error calling Greet RPC: %v", err)
	}
	fmt.Printf("Response from Greet: %v", resp.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Calling the server streaming...")
	stream, err := c.GreetManyTimes(context.Background(),
		&greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Diego"},
		})
	if err != nil {
		log.Fatalf("error calling Greet RPC: %v", err)
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
		fmt.Printf("Response from Greet Many Times: %v\n", message.GetResult())
	}
}