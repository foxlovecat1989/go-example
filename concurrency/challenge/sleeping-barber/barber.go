package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Barber struct {
	ID      int
	IsSleep bool
}

func (b *Barber) cutHair(client *Client) {
	fmt.Printf("\tbarber%d is cutting client%d's hair...\n", b.ID, client.ID)
	time.Sleep(time.Duration(rand.Intn(10)*100) * time.Millisecond)
	fmt.Printf("\tbarber%d finished cut client%d's hair...\n", b.ID, client.ID)
}
