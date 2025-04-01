package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Barber struct {
	ID   int
	Name string
}

func (b *Barber) cutHair(client *Client) {
	time.Sleep(time.Duration(rand.Intn(10)*100) * time.Millisecond)
	fmt.Printf("\tbarber%d %s finished cut client%d hair...\n", b.ID, b.Name, client.ID)
}

func (b *Barber) fallAsleep() {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	fmt.Printf("\tbarber%d %s is fallAsleep...\n", b.ID, b.Name)
}
