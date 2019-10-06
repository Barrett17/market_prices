package datagrabber

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
var wg sync.WaitGroup

// Poll new data each 60 seconds
const Quantum = 10
const Timeout = time.Minute*2

func init() {
	go func() {
		ticker := time.NewTicker(time.Second * Quantum)
		defer ticker.Stop()

		started := time.Now()
		for {
			pollData();

			now := <-ticker.C
			if now.Sub(started) > Timeout {
				log.Printf("pollData() has timed out!");
			}
		}
	}()
}

func pollData() {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Println("pollData");
}

func GetData() {
	mutex.Lock()
	defer mutex.Unlock()

	
}
