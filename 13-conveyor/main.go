package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type pair struct {
	word, reversed string
}

// генерит случайные слова из 5 букв
// с помощью randomWord(5)
func generate(cancel <-chan struct{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case out <- randomWord(5):
			case <-cancel:
				return
			}
		}
	}()
	return out
}

// выбирает слова, в которых не повторяются буквы,
// abcde - подходит
// abcda - не подходит
func takeUnique(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for w := range in {
			m := make(map[rune]struct{}, 0)
			res := true
			for _, c := range w {
				if _, ok := m[c]; ok {
					res = false
					break
				}
				m[c] = struct{}{}
			}
			if res {
				select {
				case out <- w:
				case <-cancel:
					return
				}
			}
		}
	}()
	return out
}

// переворачивает слова
// abcde -> edcba
func reverse(cancel <-chan struct{}, in <-chan string) <-chan pair {
	out := make(chan pair)
	go func() {
		defer close(out)
		for w := range in {
			r := []rune(w)
			for i := 0; i < len(r)/2; i++ {
				r[i], r[len(r)-i-1] = r[len(r)-i-1], r[i]
			}
			select {
			case out <- pair{w, string(r)}:
			case <-cancel:
				return
			}
		}
	}()
	return out
}

// объединяет c1 и c2 в общий канал
func merge(cancel <-chan struct{}, c1, c2 <-chan pair) <-chan pair {
	out := make(chan pair)
	var wg sync.WaitGroup
	wg.Add(2)

	f := func(c <-chan pair) {
		defer wg.Done()
		for v := range c {
			select {
			case out <- v:
			case <-cancel:
				return
			}
		}
	}

	go f(c1)
	go f(c2)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// печатает первые n результатов
func print(cancel <-chan struct{}, in <-chan pair, n int) {
	for i := 0; i < n; i++ {
		select {
		case p := <-in:
			fmt.Printf("%v -> %v\n", p.word, p.reversed)
		case <-cancel:
			return
		}
	}
}

// генерит случайное слово из n букв
func randomWord(n int) string {
	const letters = "aeiourtnsl"
	chars := make([]byte, n)
	for i := range chars {
		chars[i] = letters[rand.Intn(len(letters))]
	}
	return string(chars)
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	c1 := generate(cancel)
	c2 := takeUnique(cancel, c1)
	c3_1 := reverse(cancel, c2)
	c3_2 := reverse(cancel, c2)
	c4 := merge(cancel, c3_1, c3_2)
	print(cancel, c4, 10)
}
