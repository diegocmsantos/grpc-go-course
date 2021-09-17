package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

const port = ":50051"

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

func main() {
	// if we crash the code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Server Started")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting the server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for a signal to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a shutdown signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("End of Program")
}
