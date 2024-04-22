// Работяга
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Worker представляет собой обработчик, который можно запустить, а потом остановить.
type Worker struct {
	fn     func()
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewWorker создает новый экземпляр обработчика.
func NewWorker(fn func()) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{fn, ctx, cancel, sync.WaitGroup{}}
}

// Start запускает в отдельной горутине цикл, в котором
// вызывает раз за разом функцию fn.
// Гарантируется, что клиент вызывает Start только один раз.
// Гарантируется, что клиент не вызывает Start после вызова Stop.
func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.ctx.Done():
				return
			default:
				w.fn()
			}
		}
	}()
}

// Stop останавливает горутину, запущенную в Start.
// Гарантируется, что клиент не вызывает Stop до вызова Start.
// Клиент может вызвать Stop несколько раз, в том числе из разных горутин.
// После первого успешного вызова Stop повторные вызовы не приводят к ошибке.
func (w *Worker) Stop() {
	w.cancel()
}

// Wait блокирует вызвавшую его горутину до тех пор, пока обработчик не будет остановлен.
// Гарантируется, что клиент не вызывает Wait до вызова Start.
// Клиент может вызвать Wait несколько раз, в том числе из разных горутин.
// Клиент может вызвать Wait после Stop.
// Вызов Wait после Stop не приводит к блокировке.
func (w *Worker) Wait() {
	w.wg.Wait()
}

func main() {
	fn := func() {
		fmt.Print(".")
		time.Sleep(10 * time.Millisecond)
	}

	w := NewWorker(fn)
	w.Start()

	// эта горутина остановит работягу через 50 мс
	go func() {
		time.Sleep(50 * time.Millisecond)
		w.Stop()
	}()

	// подождем, пока кто-нибудь остановит работягу
	w.Wait()
	fmt.Println("done")
}
