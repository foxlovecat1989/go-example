package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pg := NewPressureGauge(10)
	wg := &sync.WaitGroup{}
	ticker := time.NewTicker(100 * time.Millisecond)

	var i int
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("shut down...\n")
			return
		case <-ticker.C:
			wg.Add(1)
			i++
			go func(i int) {
				defer wg.Done()

				fmt.Printf("start processing task%d...\n", i)

				err := pg.Process(func() {
					doThingThatShouldBeLimited()
				})
				if err != nil {
					fmt.Printf("task%d,  err: %+v", i, err)
					return
				}

				fmt.Printf("processing task%d done...\n", i)
			}(i)
		}
	}
}

func doThingThatShouldBeLimited() string {
	time.Sleep(1 * time.Second)
	return "done"
}
