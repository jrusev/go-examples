package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	iterationFlag = flag.Int("iteration", 5, "iteration to run")
	jcFlag        = flag.Bool("jc", false, "third chef?")
)

var logs = map[time.Time]string{}

type order struct {
	dish     string
	num      int
	duration time.Duration
}

func newChef(name string, station int) *chef {
	return &chef{station: station, name: name, mutex: &sync.Mutex{}}
}

type chef struct {
	station   int
	name      string
	busyState bool

	mutex *sync.Mutex
}

func (c *chef) isBusy() bool {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	return c.busyState
}

func (c *chef) busy(b bool) {
	c.mutex.Lock()
	c.busyState = b
	c.mutex.Unlock()
}

func (c *chef) cook(o *order) {
	if c.isBusy() {
		fmt.Println("hold on!! ", c.name, "is already cooking!")
		for c.isBusy() {
			<-time.After(10 * time.Millisecond)
		}
	}
	c.busy(true)
	fmt.Printf("\t%s is cooking a %s (order %d)\n", c.name, o.dish, o.num)
	time.Sleep(o.duration)
	logs[time.Now()] = fmt.Sprintf("%s cooked a %s", c.name, o.dish)
	c.rest()
}

func (c *chef) rest() {
	// minimal rest required by law
	fmt.Printf("\t\t*%s is resting*\n", c.name)
	time.Sleep(100 * time.Millisecond)
	c.busy(false)
}

func main() {
	flag.Parse()

	chefs := []*chef{newChef("Alice", 2), newChef("Bob", 3)}
	if *jcFlag {
		chefs = append(chefs, newChef("Jean-Claude", 5))
	}
	orders := []*order{
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

	startT := time.Now()
	switch *iterationFlag {
	case 1:
		one(orders, chefs)
	case 2:
		two(orders, chefs)
	case 3:
		three(orders, chefs)
	case 4:
		four(orders, chefs)
	case 5:
		five(orders, chefs)
	}

	fmt.Printf("all done in %s, closing the kitchen\n", time.Since(startT))
	fmt.Println("logs:")
	for t, entry := range logs {
		fmt.Printf("%s: %s\n", t, entry)
	}
	os.Exit(0)
}

// round robin, only 1 chef cooks everything :(
func one(orders []*order, chefs []*chef) {
	for _, order := range orders {

		for _, chef := range chefs {
			if !chef.isBusy() {
				chef.cook(order)
				break
			} else {
				fmt.Println(".")
			}
		}

	}
}

// exit right away because we aren't waiting
func two(orders []*order, chefs []*chef) {
	for _, order := range orders {

		for _, chef := range chefs {
			if !chef.isBusy() {
				go chef.cook(order)
				break
			} else {
				fmt.Println(".")
			}
		}

	}
}

/*
	use a wait group to wait for the chefs to be done
	try to send the order to a chef until one is available
	ok, but confusing code
*/
func three(orders []*order, chefs []*chef) {
	wg := &sync.WaitGroup{}
	for _, order := range orders {
	outterLoop:
		for {
			for _, chef := range chefs {
				if !chef.isBusy() {
					chef.cookAndYell(order, wg)
					break outterLoop
				}
			}
		}
	}
	wg.Wait()
}

func (c *chef) cookAndYell(o *order, wg *sync.WaitGroup) {
	wg.Add(1)
	if c.isBusy() {
		fmt.Printf("\t%s is already cooking %s (order %d)\n", c.name, o.dish, o.num)
		for c.isBusy() {
			<-time.After(10 * time.Millisecond)
		}
	}
	c.busy(true)
	go func() {
		fmt.Printf("\t%s is cooking a %s (order %d)\n", c.name, o.dish, o.num)
		time.Sleep(o.duration)
		//logs[time.Now()] = fmt.Sprintf("%s cooked a %s", c.name, o.dish)
		c.restAndYell(wg)
	}()
}

func (c *chef) restAndYell(wg *sync.WaitGroup) {
	// minimal rest required by law
	fmt.Printf("\t\t*%s is resting*\n", c.name)
	time.Sleep(500 * time.Millisecond)
	c.busy(false)
	wg.Done()
}

// almost proper solution
func four(orders []*order, chefs []*chef) {
	orderWheel := make(chan *order)
	for _, c := range chefs {
		go func(c *chef) {
			for o := range orderWheel {
				c.cook(o)
			}
		}(c)
	}

	for _, order := range orders {
		orderWheel <- order
	}
	close(orderWheel)
}

//  55555555555555555555555
func five(orders []*order, chefs []*chef) {
	wg := &sync.WaitGroup{}
	orderWheel := make(chan *order)
	for _, c := range chefs {
		wg.Add(1)
		go func(c *chef) {
			for o := range orderWheel {
				c.cook(o)
			}
			wg.Done()
		}(c)
	}

	for _, order := range orders {
		orderWheel <- order
	}
	close(orderWheel)
	wg.Wait()
}
