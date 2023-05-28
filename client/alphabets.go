package main

import (
	pb "Alphabet-Count/proto/generated"
	"context"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

func SendAlphabets(c pb.CounterClient, sid int64) {
	sTime := time.Now()
	sReqs := sync.Map{}
	//goland:noinspection SpellCheckingInspection
	strm, err := c.Alphabet(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream %v\n", err)
	}

	//goland:noinspection SpellCheckingInspection
	waitchn := make(chan struct{})
	go func() {
		for i := 1; i <= 4096; i++ {
			message := GetLetterMessage(sid, int64(i))
			sReqs.Store(message.MessageId, 1)
			log.Printf("Sending message %v\n", message)
			_ = strm.Send(message)
			//time.Sleep(1 * time.Second)
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

			val, ok := sReqs.LoadAndDelete(res.MessageId)
			if !ok {
				log.Fatalf("Received invalid response %v\n", val)
			}
			//log.Printf("Received response %v\n", res)
		}

		sReqs.Range(func(key, value any) bool {
			log.Fatalf("Missing responses for requests")
			return false
		})

		close(waitchn)
		eTime := time.Since(sTime)
		log.Printf("Spanned %s to complete 4096 request", eTime)
	}()

	<-waitchn
}

func GetLetterMessage(sid int64, num int64) *pb.LetterMessage {
	timestamp := time.Now().UnixNano()
	r := GetRandomInstance()
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