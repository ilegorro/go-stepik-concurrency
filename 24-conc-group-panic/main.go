// Concurrent-группа с паникой
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// ConcGroup выполняет присылаемую работу в отдельных горутинах.
type ConcGroup struct {
	wg       sync.WaitGroup
	panicMsg string
}

// NewConcGroup создает новый экземпляр ConcGroup.
func NewConcGroup() *ConcGroup {
	return &ConcGroup{sync.WaitGroup{}, ""}
}

// Run выполняет присланную работу в отдельной горутине.
// Если горутина запаниковала, Run не паникует.
func (p *ConcGroup) Run(work func()) {
	p.wg.Add(1)
	go func() {
		defer func() {
			val := recover()
			if val != nil {
				p.panicMsg = val.(string)
			}
			p.wg.Done()
		}()
		work()
	}()
}

// Wait ожидает, пока не закончится вся выполняемая в данный момент работа.
// Если запаниковала хотя бы одна из горутин, запущенных через Run -
// Wait тоже паникует.
func (p *ConcGroup) Wait() {
	p.wg.Wait()
	if p.panicMsg != "" {
		panic(p.panicMsg)
	}
}

func main() {
	work := func() {
		if rand.Intn(4) == 1 {
			panic("oopsie")
		}
		// do stuff
	}

	defer func() {
		val := recover()
		if val == nil {
			fmt.Println("work done")
		} else {
			fmt.Printf("panicked: %v!\n", val)
		}
	}()

	p := NewConcGroup()

	for i := 0; i < 4; i++ {
		p.Run(work)
	}

	p.Wait()
}
