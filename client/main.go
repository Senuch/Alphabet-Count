package main

import (
	pb "Alphabet-Count/proto/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const ADDRESS string = "0.0.0.0:50001"

var EXIT = make(chan bool)

func main() {
	con, err := grpc.Dial(ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect with server %v\n", err)
	}

	defer func(con *grpc.ClientConn) {
		_ = con.Close()
	}(con)

	c := pb.NewCounterClient(con)

	log.Println("Starting client")
	go SendAlphabetRequests(c)
	<-EXIT
	log.Println("Client closed")
}

func SendAlphabetRequests(c pb.CounterClient) {
	sid := int64(1)
	for {
		go SendAlphabets(c, sid)
		sid++
		time.Sleep(1 * time.Second)
	}
}