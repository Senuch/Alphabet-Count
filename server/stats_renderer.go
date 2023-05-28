package main

import (
	"log"
	"time"
)

func RenderStats() {
	for {
		ltr, _, msgCnt := COUNTER.GetCounterStats()

		if ltr == "" {
			ltr = "No Data"
		}

		log.Printf("Stats => Highest Recurring Letter: %s, Total Messages %d\n", ltr, msgCnt)
		time.Sleep(1 * time.Second)
	}
}