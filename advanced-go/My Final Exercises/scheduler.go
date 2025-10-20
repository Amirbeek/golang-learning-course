package main

import "sync"

type Scheduler struct {
	jobs    chan Transaction
	results chan string
	wg      sync.WaitGroup
}
