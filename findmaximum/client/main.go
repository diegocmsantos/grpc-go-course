package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc-go-course/findmaximum/findmaximumpb"
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

	c := findmaximumpb.NewFindMaximumServiceClient(cc)
	doBiDiStreaming(c)
}

func doBiDiStreaming(c findmaximumpb.FindMaximumServiceClient) {

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while invoking the find maximum server: %v\n", err)
	}

	numbers := []int32{1,5,3,6,2,20}

	waitCh := make(chan struct{})

	go func() {
		for _, n := range numbers {
			err := stream.Send(&findmaximumpb.FindMaximumRequest{Number: n})
			if err != nil {
				fmt.Printf("error sending message to the server: %v\n", err)
				continue
			}
			time.Sleep(1 * time.Second)
		}
		err = stream.CloseSend()
		if err != nil {
			fmt.Printf("error closing the client stream: %v\n", err)
		}
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Client has sent an end of file signal")
					break
				}
				fmt.Printf("error receiving message from the client: %v\n", err)
			}
			fmt.Printf("Max number from the client: %d\n", resp.GetAverage())
		}
		close(waitCh)
	}()

	<- waitCh
}

