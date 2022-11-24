package main

import (
	"fmt"
	"sync"
)

func merge(c ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(c))
	for _, c := range c {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func fanIn() {
	in := gen(2, 3)
	in2 := gen(4)
	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in2)

	for n := range merge(c1, c2) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}
}
