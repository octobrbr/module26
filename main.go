package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

const bufferTimer time.Duration = 5 * time.Second

const bufferSize int = 10

type RingIntBuffer struct {
	array []int
	pos   int
	size  int
	m     sync.Mutex
}

func NewRingIntBuffer(size int) *RingIntBuffer {
	return &RingIntBuffer{make([]int, size), -1, size, sync.Mutex{}}
}

func (r *RingIntBuffer) Push(el int) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.pos == r.size-1 {
		for i := 1; i <= r.size-1; i++ {
			r.array[i-1] = r.array[i]
		}
		r.array[r.pos] = el
	} else {
		r.pos++
		r.array[r.pos] = el
	}
}

func (r *RingIntBuffer) Get() []int {
	if r.pos < 0 {
		return nil
	}
	r.m.Lock()
	defer r.m.Unlock()
	var output []int = r.array[:r.pos+1]
	r.pos = -1
	return output
}

func producer(wg *sync.WaitGroup) <-chan int {
	c := make(chan int)
	go func() {
		fmt.Println("Вводите целые числа для обработки (для завершения введите stop)")
		scanner := bufio.NewScanner(os.Stdin)
		var text string
		for {
			scanner.Scan()
			text = scanner.Text()
			if text == "stop" {
				wg.Done()
				return
			}
			i, err := strconv.Atoi(text)
			if err != nil {
				fmt.Printf("'%s' не является числом!\n", text)
				continue
			}
			c <- i
		}
	}()
	return c
}

func filterNegative(source <-chan int, filtered chan<- int) {
	for i := range source {
		if i >= 0 {
			filtered <- i
		}
	}
}

func filteredNonThree(source <-chan int, filtered chan<- int) {
	for i := range source {
		if (i != 0) && (i%3 != 0) {
			filtered <- i
		}
	}
}

func bufferStage(source <-chan int, r *RingIntBuffer) {
	for i := range source {
		r.Push(i)
	}
}

func consumer(r *RingIntBuffer, t *time.Ticker) {
	for range t.C {
		b := r.Get()
		if len(b) > 0 {
			fmt.Println("Получены данные:", b)
		}
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	ch := producer(&wg)

	negativeCh := make(chan int)
	nonThreeCh := make(chan int)

	buffer := NewRingIntBuffer(bufferSize)

	go filterNegative(ch, negativeCh)

	go filteredNonThree(negativeCh, nonThreeCh)

	go bufferStage(nonThreeCh, buffer)

	go consumer(buffer, time.NewTicker(bufferTimer))

	wg.Wait()
}
