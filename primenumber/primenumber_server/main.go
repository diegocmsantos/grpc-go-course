package main

import (
	"fmt"
	"log"
	"net"

	"github.com/grpc-go-course/primenumber/primenumberpb"
	"google.golang.org/grpc"
)

const port = ":50051"

type server struct {
	primenumberpb.UnimplementedPrimeNumberServiceServer
}

func (s *server) GetPrimeNumbers(req *primenumberpb.PrimeNumberRequest, stream primenumberpb.PrimeNumberService_GetPrimeNumbersServer) error {
	k := int32(2)
	n := req.GetPrimeNumber()
	for n > 1 {
		if n % k == 0 {
			stream.SendMsg(&primenumberpb.PrimeNumberResponse{Result: k})
			n = n / k
		} else {
			k++
		}
	}
	return nil
}

func main() {
	fmt.Println("Starting server...")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	primenumberpb.RegisterPrimeNumberServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
