package main

import (
	pb "Alphabet-Count/proto/generated"
	"context"
	"io"
	"log"
	"time"
)

func sendAlphabets(c pb.CounterClient) {
	//goland:noinspection SpellCheckingInspection
	strm, err := c.Alphabet(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream %v\n", err)
	}

	reqs := []*pb.LetterMessage{
		{MessageId: 1, TimeStamp: 1, Letter: "A"},
		{MessageId: 2, TimeStamp: 2, Letter: "B"},
		{MessageId: 3, TimeStamp: 3, Letter: "C"},
	}

	//goland:noinspection SpellCheckingInspection
	waitchn := make(chan struct{})
	go func() {
		for _, req := range reqs {
			log.Printf("Sending message %v\n", req)
			_ = strm.Send(req)
			time.Sleep(1 * time.Second)
		}

		_ = strm.CloseSend()
	}()

	go func() {
		for {
			res, err := strm.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Printf("Error in receive stream %v\n", err)
				break
			}

			log.Printf("Received response %v\n", res)
		}

		close(waitchn)
	}()

	<-waitchn
}