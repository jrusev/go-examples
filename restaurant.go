package main

import (
	"fmt"
	"math/rand"
	"time"
)

func client(name string, menu []string, orders chan<- string) {
	for {
		time.Sleep(time.Duration(rand.Intn(1000))*time.Millisecond*5 + 1)
		order := menu[rand.Intn(len(menu))]
		fmt.Printf("%s ordered %s\n", name, order)
		orders <- order
	}
}

func cook(name string, orders <-chan string) {
	for {
		order := <-orders
		fmt.Printf("%s starts cooking the %s\n", name, order)
		time.Sleep(time.Duration(rand.Intn(1000))*time.Millisecond*5 + 1)
		fmt.Printf("%s cooked %s\n", name, order)
	}
}

func main() {
	menu := []string{"pizza", "pasta", "fish", "sushi", "soup"}
	orders := make(chan string)

	go client("Bob", menu, orders)
	go client("Alice", menu, orders)
	go cook("Gordon", orders)
	go cook("Jamie", orders)
	time.Sleep(time.Duration(time.Second * 5))
}
