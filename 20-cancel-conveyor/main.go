package main

import (
	"context"
	"fmt"
	"strings"
	"unicode"
)

// информация о количестве цифр в каждом слове
type counter map[string]int

// слово и количество цифр в нем
type pair struct {
	word  string
	count int
}

// считает количество цифр в словах
func countDigitsInWords(ctx context.Context, words []string) counter {
	pending := submitWords(ctx, words)
	counted := countWords(ctx, pending)
	return fillStats(ctx, counted)
}

// отправляет слова на подсчет
func submitWords(ctx context.Context, words []string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, word := range words {
			select {
			case <-ctx.Done():
				return
			case out <- word:
			}
		}
	}()
	return out
}

// считает цифры в словах
func countWords(ctx context.Context, in <-chan string) <-chan pair {
	out := make(chan pair)
	go func() {
		defer close(out)
		for {
			select {
			case word, ok := <-in:
				if !ok {
					return
				}
				count := countDigits(word)
				select {
				case <-ctx.Done():
					return
				case out <- pair{word, count}:
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// готовит итоговую статистику
func fillStats(ctx context.Context, in <-chan pair) counter {
	stats := counter{}

	for {
		select {
		case <-ctx.Done():
			return stats
		case p, ok := <-in:
			if !ok {
				return stats
			}
			stats[p.word] = p.count
		}
	}
}

// считает количество цифр в слове
func countDigits(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	words := strings.Fields(phrase)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stats := countDigitsInWords(ctx, words)
	fmt.Println(stats)
}
