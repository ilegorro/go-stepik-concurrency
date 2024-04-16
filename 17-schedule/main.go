package main

import (
	"fmt"
	"time"
)

func schedule(dur time.Duration, fn func()) func() {
	ticker := time.NewTicker(dur)
	cancel := make(chan struct{})

	cancelFunc := func() {
		select {
		case <-cancel:
		default:
			ticker.Stop()
			close(cancel)
		}
	}

	go func() {
		for {
			select {
			case <-cancel:
				return
			case <-ticker.C:
				fn()
			}
		}
	}()

	return cancelFunc
}

func main() {
	work := func() {
		at := time.Now()
		fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
	}

	cancel := schedule(50*time.Millisecond, work)
	defer cancel()

	// хватит на 5 тиков
	time.Sleep(260 * time.Millisecond)
}
