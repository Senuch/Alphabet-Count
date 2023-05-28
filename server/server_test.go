package main

import (
	pb "Alphabet-Count/proto/generated"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"testing"
)

func createServer() {
	COUNTER.Init(CHANNELBUFFERSZIE)

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

func TestMessagesSentProperly(t *testing.T) {
	go createServer()
	con, err := grpc.Dial(ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("Failed to connect with server %v\n", err)
	}

	defer func(con *grpc.ClientConn) {
		_ = con.Close()
	}(con)

	c := pb.NewCounterClient(con)
	t.Run("SendMessages", func(t *testing.T) {
		strm, err := c.Alphabet(context.Background())

		if err != nil {
			t.Fatalf("Error while creating stream %v\n", err)
		}

		reqs := []*pb.LetterMessage{
			{MessageId: 1, TimeStamp: 1, Letter: "A"},
			{MessageId: 2, TimeStamp: 2, Letter: "Z"},
			{MessageId: 3, TimeStamp: 3, Letter: "Z"},
			{MessageId: 4, TimeStamp: 4, Letter: "Z"},
		}

		waitchn := make(chan struct{})
		go func() {
			for _, req := range reqs {
				_ = strm.Send(req)
			}
			_ = strm.CloseSend()
		}()

		go func() {
			resCount := int64(0)
			ids := make(map[int64]int)
			ids[1] = 0
			ids[2] = 0
			ids[3] = 0
			ids[4] = 0
			for {
				res, err := strm.Recv()

				if err == io.EOF {
					break
				}

				if err != nil {
					t.Errorf("Error in receive stream %v\n", err)
					break
				}

				val, ok := ids[res.MessageId]
				if ok && val == 0 {
					ids[res.MessageId] = 1
				} else {
					t.Errorf("Invalid response received")
				}
				resCount += 1
			}

			letr, freq, count := COUNTER.GetCounterStats()
			if resCount == count && freq == 3 && letr == "Z" {
				log.Printf("All message received and Processed as expected")
			} else {
				log.Fatalf("Expected Letter:{Z}, Frequency:{3}, Total Letter Count:{4}. Received %s, %d, %d", letr, freq, count)
			}
			close(waitchn)
		}()

		<-waitchn
	})
}