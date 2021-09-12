package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

const port = ":50051"

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	sum := req.GetCalculator().GetA() + req.GetCalculator().GetB()
	return &calculatorpb.CalculatorResponse{Sum: sum}, nil
}

func main() {
	fmt.Println("Starting server...")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
