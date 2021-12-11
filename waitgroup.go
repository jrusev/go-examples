package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		fmt.Println("Incrementing the WaitGroup counter")
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Decrementing the WaitGroup counter")
		}()
	}
	fmt.Println("Waiting for the WaitGroup counter to go to zero")
	wg.Wait()
	fmt.Println("Done waiting")
}

// Incrementing the WaitGroup counter
// Incrementing the WaitGroup counter
// Incrementing the WaitGroup counter
// Incrementing the WaitGroup counter
// Decrementing the WaitGroup counter
// Decrementing the WaitGroup counter
// Incrementing the WaitGroup counter
// Waiting for the WaitGroup counter to go to zero
// Decrementing the WaitGroup counter
// Decrementing the WaitGroup counter
// Decrementing the WaitGroup counter
// Done waiting
