package main

import (
	"context"
	"fmt"
	"sync"
)

type Barbershop struct {
	seatCapacity int
	clientChan   chan *Client
	barbers      []*Barber
	wg           *sync.WaitGroup
}

func NewBarbershop(
	seatCapacity int,
	wg *sync.WaitGroup,
) *Barbershop {
	return &Barbershop{
		seatCapacity: seatCapacity,
		clientChan:   make(chan *Client, seatCapacity),
		barbers:      make([]*Barber, 0),
		wg:           wg,
	}
}

func (b *Barbershop) close() {
	close(b.clientChan)
}

func (b *Barbershop) addBarber(id int, name string) {
	b.barbers = append(b.barbers, &Barber{
		ID:   id,
		Name: name,
	})
}

func (b *Barbershop) Run(ctx context.Context) {
	defer b.wg.Done()

	wg := &sync.WaitGroup{}

	for _, barber := range b.barbers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case c, ok := <-b.clientChan:
					if !ok {
						fmt.Printf("the client chan is closed\n")
						return
					}

					barber.cutHair(c)
				case <-ctx.Done():
					fmt.Printf("%d:%s finished the work and go home\n", barber.ID, barber.Name)
					return
				}
			}
		}()
	}

	wg.Wait()
	fmt.Printf("all barbers finished the work...\n")
}
