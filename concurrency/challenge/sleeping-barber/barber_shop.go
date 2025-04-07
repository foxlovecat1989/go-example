package main

import (
	"fmt"
	"sync"
)

type BarberShop struct {
	clientsBufferChan chan *Client
	barbers           []*Barber
	wg                *sync.WaitGroup
}

func NewBarberShop(
	seatCapacity int,
	numOfBarbers int,
	wg *sync.WaitGroup,
) *BarberShop {
	var barbers []*Barber
	for i := 0; i < numOfBarbers; i++ {
		barbers = append(barbers, &Barber{
			id: i,
		})
	}

	return &BarberShop{
		clientsBufferChan: make(chan *Client, seatCapacity),
		barbers:           barbers,
		wg:                wg,
	}
}

func (bs *BarberShop) run() {
	defer bs.wg.Done()

	bsWg := &sync.WaitGroup{}
	for _, barber := range bs.barbers {
		bsWg.Add(1)
		go func() {
			defer bsWg.Done()

			for {
				select {
				case c, ok := <-bs.clientsBufferChan:
					if ok {
						if barber.isSleep {
							barber.wakeup(c)
						}

						barber.cutHair(c)
					} else {
						// the shop is closed
						fmt.Printf("\tbarber%d has done the today's work and left\n", barber.id)
						return
					}
				default:
					if !barber.isSleep {
						barber.sleep()
					}
				}
			}
		}()
	}

	bsWg.Wait()
	fmt.Printf("all of barbers has left...\n")
}

func (bs *BarberShop) close() {
	close(bs.clientsBufferChan)
	fmt.Printf("the shop is closed...\n")
}
