package main

import (
	"log"
	"time"
)

func RenderStats() {
	for {
		ltr, freq, msgCnt := COUNTER.GetCounterStats()
		log.Printf("Stats => Highest Recurring Letter: %s, Letter Frequency: %d, Total Messages %d\n", ltr, freq, msgCnt)
		time.Sleep(1 * time.Second)
	}
}