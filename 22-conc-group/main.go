// Concurrent-группа
package main

import (
	"fmt"
	"sync"
	"time"
)

// ConcGroup выполняет присылаемую работу в отдельных горутинах.
type ConcGroup struct {
	wg sync.WaitGroup
}

// NewConcGroup создает новый экземпляр ConcGroup.
func NewConcGroup() *ConcGroup {
	return &ConcGroup{sync.WaitGroup{}}
}

// Run выполняет присланную работу в отдельной горутине.
func (cg *ConcGroup) Run(work func()) {
	cg.wg.Add(1)
	go func() {
		work()
		cg.wg.Done()
	}()
}

// Wait ожидает, пока не закончится вся выполняемая в данный момент работа.
func (cg *ConcGroup) Wait() {
	cg.wg.Wait()
}

func main() {
	work := func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Print(".")
	}

	cg := NewConcGroup()
	for i := 0; i < 4; i++ {
		cg.Run(work)
	}
	cg.Wait()
}
