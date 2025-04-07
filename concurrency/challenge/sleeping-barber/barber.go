package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// Barber
// When a sleeping-barber has nothing to do, he or she checks the waiting room for new clients,
// and if one or more is there, a haircut takes place.
// Otherwise, the sleeping-barber goes to sleep until a new client arrives.
type Barber struct {
	id      int
	isSleep bool
}

func (b *Barber) sleep() {
	fmt.Printf("\tbarber%d is sleeping...\n", b.id)
	b.isSleep = true
}

func (b *Barber) cutHair(c *Client) {
	fmt.Printf("\tbarber%d is cutting client%d's hair\n", b.id, c.id)
	time.Sleep(time.Duration(rand.IntN(10)*100) * time.Millisecond)
	fmt.Printf("\tbarber%d finished cut client%d's hair\n", b.id, c.id)
}

func (b *Barber) wakeup(c *Client) {
	fmt.Printf("\tclint%d wake barber%d up\n", c.id, b.id)
	b.isSleep = false
}
