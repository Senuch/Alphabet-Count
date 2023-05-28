package main

import (
	"math/rand"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type Random struct {
	random *rand.Rand
}

var randomInstance *Random

func GetRandomInstance() *rand.Rand {
	if randomInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if randomInstance == nil {
			randomInstance = &Random{
				random: rand.New(rand.NewSource(time.Now().UnixNano())),
			}
		}
	}

	return randomInstance.random
}