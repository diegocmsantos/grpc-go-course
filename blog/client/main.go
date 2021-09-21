package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const port = ":50051"

func main() {
	fmt.Println("Blog client")

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

	c := blogpb.NewBlogServiceClient(cc)

	// createBlog(c)
	// readBlog(c)
	updateBlog(c)
}

func createBlog(c blogpb.BlogServiceClient) {
	blog := &blogpb.Blog{
		AuthorId: "Diego",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	fmt.Println("Creating blog")
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v\n", res.Blog)
}

func readBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Starting reading blog...")
	blogReq := &blogpb.ReadBlogRequest{BlogId: "614983dfba0ec5c19b14d5e6"}

	blogRes, err := c.ReadBlog(context.Background(), blogReq)
	if err != nil {
		grpcErr, ok := status.FromError(err)
		if ok {
			if grpcErr.Code() == codes.NotFound {
				fmt.Println(grpcErr.Message())
			}
			if grpcErr.Code() == codes.Internal {
				fmt.Printf("Unexpected error: %v\n", grpcErr.Message())
			}
			return
		}
	}

	fmt.Printf("Blog found: %v\n", blogRes.Blog)
}

func updateBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Starting updating blog...")

	blogReq := &blogpb.UpdateBlogRequest{Blog: &blogpb.Blog{
		Id:       "614849af7a0f3c911bed9cf1",
		AuthorId: "Keterin",
		Title:    "Blog Title Updated",
		Content:  "Blog content updated",
	}}

	res, err := c.UpdateBlog(context.Background(), blogReq)
	if err != nil {
		grpcErr, ok := status.FromError(err)
		if ok {
			if grpcErr.Code() == codes.NotFound {
				fmt.Println(grpcErr.Message())
			}
			if grpcErr.Code() == codes.Internal {
				fmt.Printf("Unexpected error: %v\n", grpcErr.Message())
			}
			return
		}
	}

	fmt.Printf("Blog updated: %v\n", res.Blog)
}
