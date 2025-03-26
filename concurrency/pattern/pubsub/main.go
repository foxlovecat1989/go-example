package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	producerDelayTime = 100 * time.Millisecond
	consumerDelayTime = 100 * time.Millisecond
)

type User struct {
	id int
}

func NewUser(id int) *User {
	return &User{
		id: id,
	}
}

type Producer struct {
	out   chan *User
	delay time.Duration
}

func NewProducer(out chan *User, delay time.Duration) *Producer {
	return &Producer{
		out:   out,
		delay: delay,
	}
}

func (p *Producer) produce() {
	defer close(p.out)

	for i := 0; i < 100; i++ {
		p.out <- NewUser(i)
		fmt.Printf("produce a user id: %d\n", i)
		time.Sleep(p.delay)
	}
}

type Consumer struct {
	in    chan *User
	delay time.Duration
}

func NewConsumer(in chan *User, delay time.Duration) *Consumer {
	return &Consumer{
		in:    in,
		delay: delay,
	}
}

func (c *Consumer) consume() {
	for u := range c.in {
		fmt.Printf("\t consume a user id: %d\n", u.id)
	}
}

func main() {
	userCh := make(chan *User)

	p := NewProducer(userCh, producerDelayTime)
	c := NewConsumer(userCh, consumerDelayTime)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		p.produce()
	}()

	go func() {
		defer wg.Done()
		c.consume()
	}()

	wg.Wait()
	fmt.Printf("all done")
}
