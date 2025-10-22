package main

import (
	"context"
	"sync"
	"time"
)

type Job interface {
	ID() string
	Priority() int
	Run(ctx context.Context) (any, error)
}
type PriorityQueue *[]Pool

type Pool struct {
	worker   int
	timeout  time.Duration
	queue    PriorityQueue
	canceled map[string]struct{}
	metrics  *Metrics
	mu       sync.Mutex
	cond     *sync.Cond
}

type Metrics struct {
	Success  int
	Failure  int
	Canceled int
	mu       sync.Mutex
}

func main() {

}
