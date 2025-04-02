package main

import (
	"fmt"
	"sync"
	"time"
)

type Barbershop struct {
	seatCapacity int
	clientsChan  chan *Client
	doneChan     chan struct{}
	barbers      []*Barber
	wg           *sync.WaitGroup
	isClosed     bool
}

func NewBarbershop(
	seatCapacity int,
	wg *sync.WaitGroup,
	doneChan chan struct{},
) *Barbershop {
	return &Barbershop{
		doneChan:     doneChan,
		seatCapacity: seatCapacity,
		clientsChan:  make(chan *Client, seatCapacity),
		barbers:      make([]*Barber, 0),
		wg:           wg,
		isClosed:     false,
	}
}

func (b *Barbershop) close() {
	close(b.clientsChan)
	b.isClosed = true
}

func (b *Barbershop) addBarber(id int) {
	b.barbers = append(b.barbers, &Barber{
		ID:      id,
		IsSleep: false,
	})
}

func (b *Barbershop) Run() {
	defer b.wg.Done()

	wg := &sync.WaitGroup{}
	for _, barber := range b.barbers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case client, ok := <-b.clientsChan:
					if ok {
						if barber.IsSleep {
							barber.IsSleep = false
							fmt.Printf("\tclient%d wake barber%d up\n", client.ID, barber.ID)
						}

						barber.cutHair(client)
					} else {
						if b.isClosed {
							fmt.Printf("\tbarber%d finished today's work and left\n", barber.ID)
							return
						}
					}
				default:
					barber.IsSleep = true
					fmt.Printf("\tbarber%d is snapping.\n", barber.ID)
					time.Sleep(30 * time.Millisecond)
				}
			}
		}()
	}

	wg.Wait()
}
