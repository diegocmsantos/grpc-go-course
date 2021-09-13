package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	//doServerStreaming(c)
	doClientStreaming(c)
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Calling the client streaming...")
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error calling Long Greet RPC: %v", err)
	}

	requests := []*greetpb.LongGreetRequest{
		{Greeting: &greetpb.Greeting{
			FirstName: "Diego",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Keterin",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Caio",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Paulo",
		}},
	}

	for _, greet := range requests {
		err := stream.Send(greet)
		if err != nil {
			fmt.Printf("error while sending request to the server: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("error while receiving response from the server: %v", err)
	}
	fmt.Printf("Response from Long Greet: %v", resp.Result)
}
