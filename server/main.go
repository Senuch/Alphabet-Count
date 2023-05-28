package main

import (
	pb "Alphabet-Count/proto/generated"
	"google.golang.org/grpc"
	"log"
	"net"
)

const ADDRESS string = "0.0.0.0:50001"

var COUNTER AlphabetCounter

type Server struct {
	pb.CounterServer
}

func main() {
	go RenderStats()
	lis, err := net.Listen("tcp", ADDRESS)

	if err != nil {
		log.Fatalf("Server failed to listen %v\n", err)
	}

	log.Printf("Server listening on %v\n", ADDRESS)

	s := grpc.NewServer()
	pb.RegisterCounterServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve on %v\n", err)
	}
}