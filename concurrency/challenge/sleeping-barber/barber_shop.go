package main

import (
	"fmt"
	"sync"
	"time"
)

type BarberShop struct {
	// a finite number of barbers
	barbers []*Barber
	// a finite number of seats
	clientsChan chan *Client
	// a fixed length of time the barbershop is open
	openDuration time.Duration
	wg           *sync.WaitGroup
}

func NewBarberShop(
	numOfBarbers int,
	clientsChan chan *Client,
	wg *sync.WaitGroup,
) *BarberShop {
	var barbers []*Barber
	for i := 0; i < numOfBarbers; i++ {
		barbers = append(barbers, &Barber{id: i, isSleep: false})
	}

	return &BarberShop{
		barbers:      barbers,
		clientsChan:  clientsChan,
		openDuration: openDuration,
		wg:           wg,
	}
}

func (bs *BarberShop) close() {
	close(bs.clientsChan)
	fmt.Printf("shop is closing...\n")
}

func (bs *BarberShop) addBarber(b *Barber) {
	bs.barbers = append(bs.barbers, b)
}

func (bs *BarberShop) Run() {
	defer bs.wg.Done()

	bsWg := &sync.WaitGroup{}
	for _, barber := range bs.barbers {
		bsWg.Add(1)
		go func() {
			defer bsWg.Done()

			for {
				select {
				case c, ok := <-bs.clientsChan:
					if ok {
						if barber.isSleep {
							barber.wakeup(c)
						}

						barber.cutHair(c)
					} else {
						// the shop is closed, no more clients will come
						fmt.Printf("\tbarber%d has done today's work and left\n", barber.id)
						return
					}
				default:
					if !barber.isSleep {
						fmt.Printf("\tbarber%d is going to sleep\n", barber.id)
						barber.sleep()
					}
				}
			}
		}()
	}

	bsWg.Wait()
	fmt.Printf("all of barber left...\n")
}
