package main

import (
	"fmt"
	"sync"
	"time"
)

// The Dining Philosophers problem is well known in computer science circles.
// Five philosophers, numbered from 0 through 4, live in a house where the
// table is laid for them; each philosopher has their own place at the table.
// Their only difficulty – besides those of philosophy – is that the dish
// served is a very difficult kind of spaghetti which has to be eaten with
// two forks. There are two forks next to each plate, so that presents no
// difficulty. As a consequence, however, this means that no two neighbours
// may be eating simultaneously, since there are five philosophers and five forks.

type philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

const (
	eatTime   = 2 * time.Second
	thinkTime = 4 * time.Second
	eatTimes  = 3
)

func main() {
	// print out a welcome message
	fmt.Println("Dining Philosophers Problem")
	fmt.Println("---------------------------")
	fmt.Println("The table is empty.")

	// create 5 forks
	forks := []*sync.Mutex{{}, {}, {}, {}, {}}

	// create 5 philosophers
	philosophers := []*philosopher{
		{name: "p1", leftFork: 0, rightFork: 1},
		{name: "p2", leftFork: 1, rightFork: 2},
		{name: "p3", leftFork: 2, rightFork: 3},
		{name: "p4", leftFork: 3, rightFork: 4},
		{name: "p5", leftFork: 4, rightFork: 0},
	}

	// We want everyone to be seated before they start eating, so create a WaitGroup for that, and set it to 5.
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))
	// wg is the WaitGroup that keeps track of how many philosophers are still at the table. When
	// it reaches zero, everyone is finished eating and has left. We add 5 (the number of philosophers) to this
	// wait group.
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	// dining
	for _, f := range philosophers {
		go func() {
			dine(f, forks, wg, seated)
		}()
	}

	wg.Wait()
	// print finish msg
	fmt.Printf("all philosophers have done\n")
}

func dine(philosopher *philosopher, forks []*sync.Mutex, wg *sync.WaitGroup, seated *sync.WaitGroup) {
	defer wg.Done()

	// seat the philosopher at the table
	fmt.Printf("\t%s is seated at the table.\n", philosopher.name)
	// Decrement the seated WaitGroup by one.
	seated.Done()
	// Wait until everyone is seated.
	seated.Wait()

	for i := 0; i < eatTimes; i++ {
		fmt.Printf("\t%s start to eat\n", philosopher.name)

		leftFork := philosopher.leftFork
		rightFork := philosopher.rightFork
		if leftFork > rightFork {
			leftFork, rightFork = rightFork, leftFork
		}

		forks[leftFork].Lock()
		fmt.Printf("\t%s take no.%d fork\n", philosopher.name, leftFork)
		forks[rightFork].Lock()
		fmt.Printf("\t%s take no.%d fork\n", philosopher.name, rightFork)

		fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[leftFork].Unlock()
		fmt.Printf("\t%s put down no.%d fork\n", philosopher.name, leftFork)
		forks[rightFork].Unlock()
		fmt.Printf("\t%s put down no.%d fork\n", philosopher.name, rightFork)
	}

	fmt.Printf("\t%s finished and left the seat\n", philosopher.name)
}
