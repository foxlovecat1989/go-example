// This is a simple demonstration of how to solve the Sleeping Barber dilemma, a classic computer science problem
// which illustrates the complexities that arise when there are multiple operating system processes. Here, we have
// a finite number of barbers, a finite number of seats in a waiting room, a fixed length of time the barbershop is
// open, and clients arriving at (roughly) regular intervals. When a sleeping-barber has nothing to do, he or she checks the
// waiting room for new clients, and if one or more is there, a haircut takes place. Otherwise, the sleeping-barber goes to
// sleep until a new client arrives. So the rules are as follows:
//
//   - if there are no customers, the sleeping-barber falls asleep in the chair
//   - a customer must wake the sleeping-barber if he is asleep
//   - if a customer arrives while the sleeping-barber is working, the customer leaves if all chairs are occupied and
//     sits in an empty chair if it's available
//   - when the sleeping-barber finishes a haircut, he inspects the waiting room to see if there are any waiting customers
//     and falls asleep if there are none
//   - shop can stop accepting new clients at closing time, but the barbers cannot leave until the waiting room is
//     empty
//   - after the shop is closed and there are no clients left in the waiting area, the sleeping-barber
//     goes home
//
// The Sleeping Barber was originally proposed in 1965 by computer science pioneer Edsger Dijkstra.
//
// The point of this problem, and its solution, was to make it clear that in a lot of cases, the use of
// semaphores (mutexes) is not needed.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// variables
const (
	seatCapacity = 10
	numBarber    = 4
)

func main() {
	fmt.Println("sleeping barber problem...")

	// seed our random number generator
	rand.NewSource(time.Now().UnixNano())

	// init
	wg := &sync.WaitGroup{}
	closeChan := make(chan struct{})

	// create the barbershop
	barbershop := NewBarbershop(seatCapacity, wg, closeChan)

	// add barbers
	for i := 0; i < numBarber; i++ {
		barbershop.addBarber(i)
	}
	barbershop.addBarber(1)

	// add clients
	wg.Add(1)
	go func() {
		defer wg.Done()

		var i int
		for {
			time.Sleep(time.Duration(rand.Intn(10)*50/numBarber) * time.Millisecond)

			i++
			client := &Client{ID: i}
			select {
			case <-closeChan:
				fmt.Printf("close the client channel...\n")
				barbershop.close()
				return
			case barbershop.clientsChan <- client:
				fmt.Printf("client%d enter the room\n", i)
			default:
				fmt.Printf("client%d leave the room(the room full)\n", i)
			}
		}
	}()

	// start barbershop
	wg.Add(1)
	go barbershop.Run()

	// close shop
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-time.After(10 * time.Second)
		fmt.Printf("closing the shop...\n")
		closeChan <- struct{}{}
	}()

	wg.Wait()
	fmt.Printf("main process is exiting...\n")
}
