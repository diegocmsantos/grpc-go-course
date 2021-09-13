package main

import (
	"fmt"
	"io"

	"github.com/grpc-go-course/average/averagepb"
)

type server struct {

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

		sum += req.GetNumber()
		counter++
	}
}
