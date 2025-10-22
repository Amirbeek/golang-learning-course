package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	result := make(chan int)
	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 20; i++ {
			ch <- i
		}
		close(ch)
	}()

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for num := range ch {
				square := num * num
				fmt.Printf("Consumers %d and got %d\n", id, square)
				result <- square
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(result)
	}()
	total := 0
	for num := range result {
		total += num
	}
	fmt.Printf("Total: %d\n", total)
}
