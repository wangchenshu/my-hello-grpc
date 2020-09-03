package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"hello/pb"
	"time"
)

const (
	address    = "localhost:50051"
	defaultMsg = "hello"
)

func sayHello3(c pb.GreeterClient) {
	var err error
	stream, err := c.SayHello3(context.Background())
	if err != nil {
		log.Printf("failed to call: %v", err)
		return
	}
	var i int64
	for {
		stream.Send(&pb.HelloRequest{Name: "walter"})
		if err != nil {
			log.Printf("failed to send: %v", err)
			break
		}

		reply, err := stream.Recv()
		if err != nil {
			log.Printf("failed to recv: %v", err)
			break
		}
		log.Printf("Greeting: %s", reply.Message)
		i++

		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	//msg := defaultMsg
	//if len(os.Args) > 1 {
	//	msg = os.Args[1]
	//}
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	sayHello3(c)
}
