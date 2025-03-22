package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func echo(ctx context.Context, in <-chan string, out chan<- string) {
	for {
		select {
		case msg := <-in:
			fmt.Printf("\tget a msg from channel: %s\n", msg)
			out <- fmt.Sprintf("%s%s", strings.ToUpper(msg), "!!!")
		case <-ctx.Done():
			fmt.Println("Timeout reached")
			return
		}
	}
}

func main() {
	in := make(chan string)
	out := make(chan string)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	fmt.Printf("echo start...\n")

	go echo(ctx, out, in)

	fmt.Printf("enter msg or q to quit\n")

	for {
		fmt.Printf("->")
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			return
		}

		if strings.ToLower(input) == "q" {
			return
		}

		// Check context before sending
		select {
		case out <- input:
			fmt.Printf("send a msg to channel: %s\n", input)
		case <-ctx.Done():
			fmt.Println("Timeout reached")
			return
		}

		// Check context before receiving
		select {
		case s := <-in:
			fmt.Printf("get a msg to channel: %s\n", s)
		case <-ctx.Done():
			fmt.Println("Timeout reached")
			return
		}
	}
}
