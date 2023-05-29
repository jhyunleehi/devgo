package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("data.txt")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	for {
		s, err := r.ReadString('\n') // delimeter로 '\n'을 설정하면, '\n'이 포함되어 있음.
		if err != nil {
			fmt.Print(err)
		}
		if s==""{
			return
		}
		c := strings.Split(s," ")
		fmt.Printf("%s_%s %s%s %s%s %s%s %s%s %s%s %s%s %s%s %s%s %s",c[0],c[1],c[2],c[3],c[4],c[5],c[6],c[7],c[8],c[9],c[10], c[11], c[12],c[13], c[14],c[15], c[16],c[17], c[18])
	}
}
