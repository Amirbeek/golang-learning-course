package main

import (
	"fmt"
	"sync"
)

type Fruits struct {
	name string
}

func processFoot(f Fruits) error {
	fmt.Printf("process foot: %s\n", f.name)
	return nil
}
func main() {
	var wg sync.WaitGroup
	fruits := []Fruits{
		{"Apple"},
		{"Orange"},
		{"Pear"},
	}
	for _, fruit := range fruits {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := processFoot(fruit)
			if err != nil {
				return
			}
		}()
	}
	wg.Wait()
	fmt.Println("All Fruits processed")
}
