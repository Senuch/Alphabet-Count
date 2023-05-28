package main

import (
	pb "Alphabet-Count/proto/generated"
	"context"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func sendAlphabets(c pb.CounterClient, sid int64) {
	//goland:noinspection SpellCheckingInspection
	strm, err := c.Alphabet(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream %v\n", err)
	}

	//goland:noinspection SpellCheckingInspection
	waitchn := make(chan struct{})
	go func() {
		for i := 1; i < 10; i++ {
			message := GetLetterMessage(sid, int64(i))
			log.Printf("Sending message %v\n", message)
			_ = strm.Send(message)
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

func GetLetterMessage(sid int64, num int64) *pb.LetterMessage {
	timestamp := time.Now().UnixNano()
	r := rand.New(rand.NewSource(timestamp))
	ascMin := 65
	ascMax := 91

	return &pb.LetterMessage{
		TimeStamp: timestamp,
		MessageId: GetLetterMessageId(sid, num),
		Letter:    string(rune(r.Intn(ascMax-ascMin) + ascMin)),
	}
}

func GetLetterMessageId(sid int64, num int64) int64 {
	concat := strconv.FormatInt(sid, 10) + strconv.FormatInt(num, 10)
	id, _ := strconv.ParseInt(concat, 10, 64)
	return id
}