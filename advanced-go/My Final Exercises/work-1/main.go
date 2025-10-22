package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job interface {
	Run(ctx context.Context)
}

type Runner struct {
	ID      int
	Timeout time.Duration
}

func (r Runner) Run(ctx context.Context) string {
	defer ctx.Done()

	select {
	case <-ctx.Done():
		return fmt.Sprintf("Error: %v")
	case <-time.After(r.Timeout):
		return fmt.Sprintf("done: %v", r.ID)
	}
}

func main() {
	var wg sync.WaitGroup
	defer wg.Wait()
	root := context.Background()
	deadline := time.Now().Add(500 * time.Millisecond)
	ctx, cancel := context.WithDeadline(root, deadline)
	defer cancel()

	var workers []Runner
	for i := 1; i <= 5; i++ {
		workerTimeout := rand.Intn(6) + 1
		workers = append(workers, Runner{
			ID:      i,
			Timeout: time.Duration(workerTimeout*100) * time.Millisecond,
		})
	}
	fmt.Println("Deadline:", deadline.Format(time.RFC3339Nano))

	result := make(chan string, len(workers))

	for _, w := range workers {
		wg.Add(1)
		go func(w Runner) {
			defer wg.Done()
			result <- w.Run(ctx)
		}(w)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	for res := range result {
		fmt.Println(res)
	}

}
