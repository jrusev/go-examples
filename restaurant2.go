package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type order struct {
	dish     string
	duration time.Duration
}

var menu = []order{
	{"Blanquette de veau", 1500 * time.Millisecond},
	{"Soupe à l'oignon", 850 * time.Millisecond},
	{"Quiche", 550 * time.Millisecond},
	{"Boeuf bourguignon", 3000 * time.Millisecond},
	{"Salade niçoise", 800 * time.Millisecond},
	{"Confit de canard", 1900 * time.Millisecond},
	{"Tarte tatin", 3200 * time.Millisecond},
	{"Ratatouille", 850 * time.Millisecond},
	{"Croque-monsieur", 700 * time.Millisecond},
	{"Jambon-beurre", 700 * time.Millisecond},
}

func client(name string, menuItems []int, ordersChan chan<- order, done chan<- string) {
	for _, i := range menuItems {
		order := menu[i]
		fmt.Printf("%s ordered %s\n", name, order.dish)
		ordersChan <- order
		time.Sleep(time.Duration(rand.Intn(1000))*time.Millisecond*5 + 1)
	}
	done <- name
}

func cook(name string, ordersChan <-chan order, wg *sync.WaitGroup) {
	for order := range ordersChan {
		fmt.Printf("%s starts cooking the %s\n", name, order.dish)
		time.Sleep(order.duration)
		fmt.Printf("%s cooked %s\n", name, order.dish)
	}
	wg.Done()
}

func main() {
	ordersChan := make(chan order)
	doneChan := make(chan string)

	go client("Bob", []int{0, 1, 2}, ordersChan, doneChan)
	go client("Alice", []int{3, 4, 5}, ordersChan, doneChan)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go cook("Gordon", ordersChan, wg)
	go cook("Jamie", ordersChan, wg)

	for c := 0; c < 2; c++ {
		fmt.Printf("%s done with orders\n", <-doneChan)
	}
	close(ordersChan)
	fmt.Println("No more orders accepted.")
	wg.Wait() // Wait for cooks to finish cooking.
}
