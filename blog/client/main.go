package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
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
