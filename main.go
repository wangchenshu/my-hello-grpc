package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"hello/pb"
	"io"
	"log"
	"net"
	"sync"
)

type server struct {
	sync.RWMutex
	cntOfMember uint32
}

func (s *server) SayHello3(gs pb.Greeter_SayHello3Server) error {
	s.Lock()
	s.cntOfMember++
	s.Unlock()

	for {
		in, err := gs.Recv()
		if err == io.EOF {
			s.Lock()
			s.cntOfMember--
			s.Unlock()
			return nil
		}

		if err != nil {
			log.Printf("failed to recv: %v", err)
			s.Lock()
			s.cntOfMember--
			s.Unlock()
			return err
		}

		log.Println("send msg hello to member", in.Name, "cnt of member:", s.cntOfMember)
		gs.Send(&pb.HelloReply{Message: "Hello " + in.Name})
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
