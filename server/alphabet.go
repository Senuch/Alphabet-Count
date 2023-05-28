package main

import (
	pb "Alphabet-Count/proto/generated"
	"io"
	"log"
)

func (s *Server) Alphabet(stream pb.Counter_AlphabetServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error reading client stream %v\n", err)
		}

		log.Printf("Alphabet received %s\n", req.Letter)
		COUNTER.AddLetter(req.Letter)

		a, b := COUNTER.GetHighestCountLetter()

		log.Printf("Current highest letter %s and count %d\n", a, b)

		res := &pb.LetterMessage{
			MessageId: req.MessageId,
		}

		err = stream.Send(res)

		if err != nil {
			log.Fatalf("Error sending response %v\n", err)
		}
	}
}