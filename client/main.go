package main

import (
	pb "Alphabet-Count/proto/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const ADDRESS string = "0.0.0.0:50001"

func main() {
	con, err := grpc.Dial(ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect with server %v\n", err)
	}

	defer func(con *grpc.ClientConn) {
		_ = con.Close()
	}(con)

	c := pb.NewCounterClient(con)

	sendAlphabets(c)
}