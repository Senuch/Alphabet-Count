package main

import (
	pb "Alphabet-Count/proto/generated"
	"google.golang.org/grpc"
	"log"
	"net"
)

var address = "0.0.0.0:50001"

type Server struct {
	pb.CounterServer
}

func main() {
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Server failed to listen %v\n", err)
	}

	log.Printf("Server listening on %v\n", address)

	s := grpc.NewServer()
	pb.RegisterCounterServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve on %v\n", err)
	}
}