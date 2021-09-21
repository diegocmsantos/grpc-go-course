package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-go-course/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const port = ":50051"

var collection *mongo.Collection

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string            `bson:"author_id"`
	Content  string            `bson:"content"`
	Title    string            `bson:"title"`
}

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

func (s *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()

	data := BlogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error creating blog: %v\n", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, fmt.Sprintf("cannot convert to objectID: %v\n", err))
	}

	return &blogpb.CreateBlogResponse{Blog: &blogpb.Blog{
		Id:       oid.Hex(),
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}}, nil
}

func (s *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error converting the blog ID: %v\n", err))
	}

	var blog BlogItem
	filter := bson.D{{"_id", oid}}

	err = collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("blog id [%s] not found", oid.Hex()))
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.AuthorID,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("Updating a blog...")
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error converting the blog ID: %v\n", err))
	}

	filter := bson.D{{"_id", oid}}
	update := bson.D{
		{"author_id", blog.GetAuthorId()},
		{"title", blog.GetTitle()},
		{"content", blog.GetContent()},
	}
	_, err = collection.ReplaceOne(ctx, filter, update)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error updating blog with ID [%s]: %v", blog.GetId(), err))
	}
	
	return &blogpb.UpdateBlogResponse{Blog: blog}, nil
}

func (s *server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("Deleting blog server side")

	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error converting the blog ID: %v\n", err))
	}

	filter := bson.D{{"_id", oid}}
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("blog id [%s] not found", oid.Hex()))
	}

	return &blogpb.DeleteBlogResponse{BlogId: oid.Hex()}, nil
}

func main() {
	// if we crash the code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Server Started")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection = client.Database("mydb").Collection("blog")

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
	fmt.Println("Closing mongodb connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of Program")
}
