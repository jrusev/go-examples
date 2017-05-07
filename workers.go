package main

import "fmt"
import "time"

func f(in <-chan int, out chan<- int) {
	for j := range in {
		time.Sleep(time.Second)
		out <- j * 2
	}
}

func main() {
	in := make(chan int, 100)
	out := make(chan int, 100)
	go f(in, out)
	for j := 1; j <= 8; j++ {
		in <- j
	}
	for a := 1; a <= 8; a++ {
		fmt.Println(<-out)
	}
}