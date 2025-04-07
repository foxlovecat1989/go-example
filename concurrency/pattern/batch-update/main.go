package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	batchSize       = 100
	tickerInterval  = 1000 * time.Millisecond
	processDuration = 10 * time.Second
	produceRate     = 5
)

func main() {
	// init
	wg := &sync.WaitGroup{}
	consumer := NewConsumer(wg, batchSize, tickerInterval)
	doneChan := make(chan struct{})
	rand.NewSource(time.Now().UnixNano())

	// producer
	var i int
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Duration(rand.Intn(5)*produceRate) * time.Millisecond)
			i++
			select {
			case <-doneChan:
				consumer.close()
				return
			case consumer.usersBufferChan <- &User{id: i}:
			}
		}
	}()

	// consumer
	wg.Add(1)
	go consumer.run()

	// stop
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-time.After(processDuration)
		fmt.Printf("shutdown the process...\n")
		// shutdown
		doneChan <- struct{}{}
	}()

	wg.Wait()
	fmt.Printf("main process exit...\n")
}
