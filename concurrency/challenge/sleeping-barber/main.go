// This is a simple demonstration of how to solve the Sleeping Barber dilemma, a classic computer science problem
// which illustrates the complexities that arise when there are multiple operating system processes. Here, we have
// a finite number of barbers, a finite number of seats in a waiting room, a fixed length of time the barbershop is
// openDuration, and clients arriving at (roughly) regular intervals. When a sleeping-barber has nothing to do, he or she checks the
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

const (
	seatCapacity = 2
	openDuration = 5 * time.Second
	numOfBarber  = 1
)

func main() {
	fmt.Printf("start sleeping barbers problem\n")
	// init
	rand.NewSource(time.Now().UnixNano())
	clientsChan := make(chan *Client, seatCapacity)
	closeChan := make(chan struct{})
	wg := &sync.WaitGroup{}
	shop := NewBarberShop(numOfBarber, clientsChan, wg)
	// producer
	wg.Add(1)
	go func() {
		defer wg.Done()

		i := 0
		for {
			i++
			// clients arriving at (roughly) regular intervals
			time.Sleep(time.Duration(rand.Intn(10)*50/numOfBarber) * time.Millisecond)
			select {
			case clientsChan <- &Client{id: i}:
				fmt.Printf("client%d is enter to the shop\n", i)
			case <-closeChan:
				shop.close()
				return
			default:
				fmt.Printf("client%d is leaveing, no room available\n", i)
			}
		}
	}()

	// consumer
	wg.Add(1)
	go shop.Run()

	// close
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-time.After(openDuration)
		closeChan <- struct{}{}
	}()

	wg.Wait()
	fmt.Printf("the main process exit...")
}
