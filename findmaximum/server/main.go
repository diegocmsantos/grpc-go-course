package main

import (
	"fmt"
	"io"

	"github.com/grpc-go-course/findmaximum/findmaximumpb"
)

type server struct {

}

func (s *server) FindMaximum(stream findmaximumpb.FindMaximumService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked by the client")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client has sent a 'End of file' signal")
				return nil
			}
			return fmt.Errorf("error while getting a message from the client: %v\n", err)
		}

		n := req.GetNumber()
		fmt.Printf("Received [%d] from the client, the current max is [%d]\n", n, max)
		if n > max {
			max = n
			fmt.Printf("new max number is [%d]\n", max)
		}
		err = stream.Send(&findmaximumpb.FindMaximumResponse{Average: max})
		if err != nil {
			return fmt.Errorf("error while responding the client: %v\n", err)
		}
	}
}

func main() {

}
