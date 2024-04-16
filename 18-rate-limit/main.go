package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrCanceled error = errors.New("canceled")

func withRateLimit(limit int, fn func()) (handle func() error, cancel func()) {
	dur := time.Duration(1000/limit) * time.Millisecond
	canceled := make(chan struct{})
	queue := make(chan int, limit)
	for i := 0; i < limit; i++ {
		queue <- i
	}

	handle = func() error {
		select {
		case <-canceled:
			return ErrCanceled
		default:
			id := <-queue
			time.Sleep(dur)
			go func() {
				fn()
				queue <- id
			}()
		}
		return nil
	}

	cancel = func() {
		select {
		case <-canceled:
		default:
			close(canceled)
		}
	}

	return
}

func main() {
	work := func() {
		fmt.Print(".")
	}

	handle, cancel := withRateLimit(5, work)
	defer cancel()

	start := time.Now()
	const n = 10
	for i := 0; i < n; i++ {
		handle()
	}
	fmt.Println()
	fmt.Printf("%d queries took %v\n", n, time.Since(start))
}
