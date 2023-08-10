package main

import (
	"fmt"

	"sync"
	"time"
)

func main1() {
	var mutex = &sync.Mutex{}

	total := 0
	for r := 0; r < 10; r++ {
		fmt.Printf("<%d>\n", r)
		go func() {
			for j := 1; j < 1000; j++ {
				fmt.Printf("%d--", j)
				mutex.Lock()
				if total < 5 {
					fmt.Printf("[%d]\n", total)
					total = total + 1
				}
				mutex.Unlock()

				if total <= 5 {
					fmt.Printf(">>>>>>>>call>>>>>> %d\n", total)
					mutex.Lock()
					fmt.Printf("--- %d\n", total)
					total = total - 1
					mutex.Unlock()
				}
				fmt.Printf("--%d\n", j)
			}
		}()

	}

	time.Sleep(time.Second * 10)

}
