package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type Consumer struct {
	wg              *sync.WaitGroup
	mx              *sync.Mutex
	usersBufferChan chan *User
	batchSize       int
	usersBuffer     []*User
	ticker          *time.Ticker
	tickerInterval  int
}

func NewConsumer(
	wg *sync.WaitGroup,
	batchSize int,
	tickerInterval time.Duration,
) *Consumer {
	return &Consumer{
		wg:              wg,
		mx:              &sync.Mutex{},
		usersBufferChan: make(chan *User, batchSize),
		usersBuffer:     make([]*User, 0, batchSize),
		ticker:          time.NewTicker(tickerInterval),
		batchSize:       batchSize,
	}
}

func (c *Consumer) run() {
	defer c.wg.Done()

	dbWg := &sync.WaitGroup{}

	for {
		select {
		case u, ok := <-c.usersBufferChan:
			if !ok {
				// update remaining
				usersCopy := c.copy()
				c.flush(dbWg, usersCopy)
				c.resetBuffer()
				fmt.Printf("\tthe channel is closed...\n")
				dbWg.Wait()
				fmt.Printf("\tall of db savinbg has done...\n")
				return
			}

			// append
			c.append(u)
			// check is up to limit
			if c.isUpToLimit() {
				usersCopy := c.copy()
				c.flush(dbWg, usersCopy)
				c.resetBuffer()
				c.ticker.Reset(tickerInterval)
			}
		case <-c.ticker.C:
			usersCopy := c.copy()
			c.flush(dbWg, usersCopy)
			c.resetBuffer()
		}
	}
}

func (c *Consumer) copy() []*User {
	c.mx.Lock()
	defer c.mx.Unlock()

	usersCopy := make([]*User, len(c.usersBuffer))
	copy(usersCopy, c.usersBuffer)

	return usersCopy
}

func (c *Consumer) flush(wg *sync.WaitGroup, users []*User) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		batchUpdate(users)
	}()
}

func (c *Consumer) append(u *User) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.usersBuffer = append(c.usersBuffer, u)
}

func (c *Consumer) isUpToLimit() bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	b := len(c.usersBuffer) == c.batchSize
	return b
}

func (c *Consumer) resetBuffer() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.usersBuffer = make([]*User, 0, c.batchSize)
}

func (c *Consumer) close() {
	fmt.Printf("\tclosing the channel...\n")
	close(c.usersBufferChan)
}

func batchUpdate(users []*User) {
	time.Sleep(time.Duration(rand.IntN(10)*10) * time.Millisecond)
	fmt.Printf("saving %d users to db ...\n", len(users))
	for _, u := range users {
		fmt.Printf("user%d\t", u.id)
	}
	fmt.Printf("\n")
}
