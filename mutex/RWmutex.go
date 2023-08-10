package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedRsc = make(map[string]string)
var mu sync.RWMutex

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	// Writer routine
	go func() {
		defer wg.Done()
		mu.RLock()
		defer mu.RUnlock()
		time.Sleep(5 * time.Second)
		sharedRsc["key1"] = "value1"
	}()

	time.Sleep(1 * time.Second)

	// Reader routine
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			mu.RLock()
			defer mu.RUnlock()

			fmt.Printf("Routine %d sees the value: %s\n", id, sharedRsc["key1"])
		}(i)
	}

	wg.Wait()
}