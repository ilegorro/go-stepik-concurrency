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

// submitWords начитывает в канал слова из генератора
func submitWords(next nextFunc, pending chan<- string) {
	for w := next(); ; w = next() {
		pending <- w
		if w == "" {
			return
		}
	}
}

// countWords получает из канала слова и передает в другой канал пары слово-количество цифр
func countWords(pending <-chan string, counted chan<- pair) {
	for w := <-pending; ; w = <-pending {
		counted <- pair{w, countDigits(w)}
		if w == "" {
			return
		}
	}
}

// fillStats заполняет карту с количеством цифр в словах
func fillStats(counted <-chan pair) counter {
	stats := counter{}
	for p := <-counted; p.word != ""; p = <-counted {
		stats[p.word] = p.count
	}

	return stats
}

// countDigitsInWords считает количество цифр в словах,
// выбирая очередные слова с помощью next()
func countDigitsInWords(next nextFunc) counter {
	pending := make(chan string)
	go submitWords(next, pending)

	counted := make(chan pair)
	go countWords(pending, counted)

	return fillStats(counted)
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
	phrase := "0ne 1wo thr33 4068"
	next := wordGenerator(phrase)
	stats := countDigitsInWords(next)
	printStats(stats)
}
