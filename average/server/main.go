package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/grpc-go-course/average/averagepb"
	"google.golang.org/grpc"
)

const port = ":50051"

type server struct {
	averagepb.UnimplementedAverageServiceServer
}

func (s *server) Average(resp averagepb.AverageService_AverageServer) error {
	sum := float32(0)
	counter := float32(0)
	for {
		req, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of file")
				return resp.SendAndClose(&averagepb.AverageResponse{Average: sum / counter})
			}
			return err
		}

		fmt.Printf("number from the client: %.2f\n", req.GetNumber())
		sum += req.GetNumber()
		counter++
	}
}

func main() {
	fmt.Println("Starting server...")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	averagepb.RegisterAverageServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
