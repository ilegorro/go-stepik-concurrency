package main

import (
	"fmt"
	"strings"
	"unicode"
)

// nextFunc возвращает следующее слово из генератора
type nextFunc func() string

// counter хранит количество цифр в каждом слове.
// ключ карты - слово, а значение - количество цифр в слове.
type counter map[string]int

// pair хранит слово и количество цифр в нем
type pair struct {
	word  string
	count int
}

// countDigitsInWords считает количество цифр в словах,
// выбирая очередные слова с помощью next()
func countDigitsInWords(next nextFunc) counter {
	pending := make(chan string)
	go submitWords(next, pending)

	done := make(chan struct{})
	counted := make(chan pair)

	for i := 0; i < 4; i++ {
		go countWords(done, pending, counted)
	}
	go func() {
		for i := 0; i < 4; i++ {
			<-done
		}
		close(counted)
	}()

	return fillStats(counted)
}

// submitWords отправляет слова на подсчет
func submitWords(next nextFunc, out chan<- string) {
	for {
		word := next()
		if word == "" {
			break
		}
		out <- word
	}
	close(out)
}

// countWords считает цифры в словах
func countWords(done chan<- struct{}, in <-chan string, out chan<- pair) {
	for word := range in {
		out <- pair{word, countDigits(word)}
	}
	done <- struct{}{}
}

// fillStats готовит итоговую статистику
func fillStats(in <-chan pair) counter {
	stats := counter{}
	for p := range in {
		stats[p.word] = p.count
	}
	return stats
}

// countDigits возвращает количество цифр в строке
func countDigits(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

// printStats печатает слова и количество цифр в каждом
func printStats(stats counter) {
	for word, count := range stats {
		fmt.Printf("%s: %d\n", word, count)
	}
}

// wordGenerator возвращает генератор, который выдает слова из фразы
func wordGenerator(phrase string) nextFunc {
	words := strings.Fields(phrase)
	idx := 0
	return func() string {
		if idx == len(words) {
			return ""
		}
		word := words[idx]
		idx++
		return word
	}
}

func main() {
	phrase := "1 22 333 4444 55555 666666 7777777 88888888"
	next := wordGenerator(phrase)
	stats := countDigitsInWords(next)
	printStats(stats)
}
