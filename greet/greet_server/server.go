package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

const port = ":50051"

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	return &greetpb.GreetResponse{
		Result: fmt.Sprintf("Hello %s", firstName),
	}, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) (error) {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{Result: fmt.Sprintf("Hello %s %d time(s)", firstName, i)}
		err := stream.SendMsg(res)
		if err != nil {
			fmt.Printf("error sending message: %v", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := "Hello "
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of file")
				fmt.Println("Sending response")
				return stream.SendAndClose(&greetpb.LongGreetResponse{Result: result})
			}
			fmt.Printf("error while receiving messages: %v", err)
		}

		firstName := resp.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
}

func main() {
	fmt.Println("Starting server...")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
