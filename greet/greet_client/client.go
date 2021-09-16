package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)
	//doServerStreaming(c)
	// doClientStreaming(c)
	//doBiDiStreaming(c)
	doUnaryWithDeadline(c, 5 * time.Second)
	// doUnaryWithDeadline(c, 1 * time.Second)
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

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Calling the BiDi streaming...")

	// create a stream by invoking the client
	cli, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while sending message to the client: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{Greeting: &greetpb.Greeting{FirstName: "Diego"}},
		{Greeting: &greetpb.Greeting{FirstName: "Keterin"}},
		{Greeting: &greetpb.Greeting{FirstName: "Caio"}},
		{Greeting: &greetpb.Greeting{FirstName: "Paulo"}},
		{Greeting: &greetpb.Greeting{FirstName: "Paulo"}},
	}

	waitChannel := make(chan struct{})

	// send a bunch of messages to the client (go routines)
	go func() {
		for _, req := range requests {
			firstName := req.GetGreeting().GetFirstName()
			fmt.Printf("sending [%s] to the client\n", firstName)
			err := cli.Send(req)
			if err != nil {
				fmt.Printf("error while sending message to the client: %v\n", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		err := cli.CloseSend()
		if err != nil {
			fmt.Printf("error while closing client stream: %v", err)
		}
	}()

	//receive a bunch of messages from the client (go routines)
	go func() {
		for {
			resp, err := cli.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("End of file")
					break
				}
				fmt.Printf("error while receiving: %v\n", err)
				break
			}
			fmt.Printf("Received from the server: %v\n", resp.GetResult())
		}
		close(waitChannel)
	}()

	// block until everything is done
	<-waitChannel
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a Unary with deadline RPC...")

	req := &greetpb.GreetWithDeadlineRequest{Greeting: &greetpb.Greeting{FirstName: "Diego", LastName: "Maia"}}

	resp, err := c.GreetWithDeadline(context.Background(), req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Deadline error: %v\n", err)
			} else {
				fmt.Printf("unexpected error: %v, code: %d\n", err, statusErr.Code())
			}
		} else {
			log.Fatalf("error receiving a response from server: %v\n", err)
		}
		return
	}
	fmt.Printf("Greet With Deadline response: %s\n", resp.GetResult())
}
