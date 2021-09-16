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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		fmt.Printf("server has received [%s] from the client\n", firstName)
		result += fmt.Sprintf("%s!", firstName)
	}
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked with a streaming request")
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of file")
				return nil
			}
			log.Fatalf("error while reading client stream: %v\n", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		fmt.Printf("Received from the client: %s\n", firstName)
		result := fmt.Sprintf("Hello %s", firstName)
		err = stream.Send(&greetpb.GreetEveryoneResponse{Result: result})
		if err != nil {
			log.Fatalf("error while sending data to the client: %v\n", err)
		}
	}
}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Println("GreetWithDeadline function was invoked with a streaming request")

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("client has canceled the request")
			return nil, status.Error(codes.DeadlineExceeded, "the client has canceled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	fmt.Printf("Received a message from the client: %s\n", firstName)
	return &greetpb.GreetWithDeadlineResponse{Result: fmt.Sprintf("Hello %s", firstName)}, nil
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
