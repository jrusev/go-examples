package main

import (
	"fmt"
	"sync"
	"time"
)

type order struct {
	dish     string
	num      int
	duration time.Duration
}

type chef string

func (c chef) cook(o order) {
	fmt.Printf("%s is cooking a %s (order %d)\n", c, o.dish, o.num)
	time.Sleep(o.duration)
	fmt.Printf("%s cooked %s (order %d)\n", c, o.dish, o.num)
}

func main() {
	chefs := []chef{"Alice", "Bob"}
	orders := []order{
		{"Blanquette de veau", 0, 1500 * time.Millisecond},
		{"Soupe à l'oignon", 1, 850 * time.Millisecond},
		{"Quiche", 2, 550 * time.Millisecond},
		{"Boeuf bourguignon", 3, 3000 * time.Millisecond},
		{"Salade niçoise", 4, 800 * time.Millisecond},
		{"Confit de canard", 5, 1900 * time.Millisecond},
		{"Tarte tatin", 6, 3200 * time.Millisecond},
		{"Ratatouille", 7, 850 * time.Millisecond},
		{"Croque-monsieur", 8, 700 * time.Millisecond},
		{"Jambon-beurre", 9, 700 * time.Millisecond},
	}

	start := time.Now()
	run(orders, chefs)
	fmt.Printf("All done in %s, closing the kitchen\n", time.Since(start))
}

func run(orders []order, chefs []chef) {
	wg := &sync.WaitGroup{}
	orderChan := make(chan order)
	for _, c := range chefs {
		wg.Add(1)
		go func(c chef) {
			for o := range orderChan {
				c.cook(o)
			}
			wg.Done()
		}(c)
	}

	for _, order := range orders {
		orderChan <- order
	}
	close(orderChan)
	wg.Wait()
}
