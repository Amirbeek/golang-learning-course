package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Student struct {
	name string
}

func TakeExam(ctx context.Context, s Student) error {
	select {
	case <-time.After(4 * time.Second):
		fmt.Printf("%s completed examination\n", s.name)
	case <-ctx.Done():
		fmt.Printf("%s canceled exam: %v\n", s.name, ctx.Err())
		return ctx.Err()
	}
	return nil
}

func main() {
	var wg sync.WaitGroup
	defer wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	students := []Student{
		{name: "Anvar"},
		{name: "Diyor"},
		{name: "Rustam"},
	}

	for _, student := range students {
		wg.Add(1)
		fmt.Println("Starting exam for:", student.name)
		go func(s Student) {
			defer wg.Done()
			err := TakeExam(ctx, s)
			if err != nil {
				fmt.Printf("%s encountered error: %v\n", s.name, err)
			}
		}(student)
	}

	wg.Wait()
	fmt.Println("All students finished or exam time ended")

}
