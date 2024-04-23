// Ограничитель вызовов
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrBusy     = errors.New("busy")
	ErrCanceled = errors.New("canceled")
)

type RateLimit struct {
	mu    sync.RWMutex
	count int
}

// withRateLimit следит, чтобы функция fn выполнялась не более limit раз в секунду.
// Возвращает функции handle (выполняет fn с учетом лимита) и cancel (останавливает ограничитель).
func withRateLimit(limit int, fn func()) (handle func() error, cancel func()) {
	r := RateLimit{}
	done := make(chan struct{})

	t := time.Tick(1 * time.Second)
	go func() {
		select {
		case <-done:
			t = nil
			return
		case <-t:
			r.mu.Lock()
			r.count = 0
			r.mu.Unlock()
		}
	}()

	handle = func() error {
		select {
		case <-done:
			return ErrCanceled
		default:
			r.mu.RLock()
			if r.count >= limit {
				r.mu.RUnlock()
				return ErrBusy
			}
			r.mu.RUnlock()
			r.mu.Lock()
			r.count++
			r.mu.Unlock()
			fn()
		}

		return nil
	}
	cancel = func() {
		select {
		case <-done:
			return
		default:
			close(done)
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

	const n = 8
	var nOK, nErr int
	for i := 0; i < n; i++ {
		err := handle()
		if err == nil {
			nOK += 1
		} else {
			nErr += 1
		}
	}
	fmt.Println()
	fmt.Printf("%d calls: %d OK, %d busy\n", n, nOK, nErr)
}
