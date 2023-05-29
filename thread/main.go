package main

import (
	"context"
	"fmt"
	"time"
)

var key bool

func main() {
	key = true
	for {
		fmt.Printf("=======>>call go()\n")
		if key {
			fmt.Printf("---->in key\n")			
			go Find()
		}
		fmt.Printf("=======>>call END\n")
		time.Sleep(100 * time.Millisecond)
	}
	
}

type Response struct {
	data   interface{}
	status bool
}

func Find() {
	key = false
	fmt.Printf("=======>>Find() call\n")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan Response, 1)

	go func() {
		time.Sleep(1 * time.Second)

		select {
		default:
			ch <- Response{data: "data", status: true}
		case <-ctx.Done():
			fmt.Printf("Canceled by timeout")
			return
		}
	}()

	select {
	case r := <-ch:
		fmt.Printf("Read from ch[%+v]\n", r)
	case <-time.After(1500 * time.Millisecond):
		fmt.Printf("Timed out")
	}
	key = true
}
