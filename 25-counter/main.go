// Безопасная карта
package main

import (
	"fmt"
	"sync"
)

// Counter представляет безопасную карту частот слов.
// Ключ - строка, значение - целое число.
type Counter struct {
	mu sync.Mutex
	m  map[string]int
}

// Increment увеличивает значение по ключу на 1.
func (c *Counter) Increment(str string) {
	c.mu.Lock()
	c.m[str]++
	c.mu.Unlock()
}

// Value возвращает значение по ключу.
func (c *Counter) Value(str string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.m[str]
}

// Range проходит по всем записям карты,
// и для каждой вызывает функцию fn, передавая в нее ключ и значение.
func (c *Counter) Range(fn func(key string, val int)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.m {
		fn(k, v)
	}
}

// NewCounter создает новую карту частот.
func NewCounter() *Counter {
	return &Counter{sync.Mutex{}, make(map[string]int)}
}

func main() {
	counter := NewCounter()

	var wg sync.WaitGroup
	wg.Add(3)

	increment := func(key string, val int) {
		defer wg.Done()
		for ; val > 0; val-- {
			counter.Increment(key)
		}
	}

	go increment("one", 100)
	go increment("two", 200)
	go increment("three", 300)

	wg.Wait()

	fmt.Println("two:", counter.Value("two"))

	fmt.Print("{ ")
	counter.Range(func(key string, val int) {
		fmt.Printf("%s:%d ", key, val)
	})
	fmt.Println("}")
}
