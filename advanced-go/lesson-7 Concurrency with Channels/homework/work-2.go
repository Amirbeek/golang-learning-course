package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type WebUrl struct {
	Url string
}

func GetResponse(ctx context.Context, url WebUrl) string {
	randomNum := rand.Intn(10) + 1 // 1–2 sec delay

	select {
	case <-ctx.Done():
		fmt.Printf("Context canceled: %v, URL: %v (delay: %d sec)\n", ctx.Err(), url.Url, randomNum)
	case <-time.After(time.Duration(randomNum) * time.Second):
		fmt.Printf("Completed request to %v (delay: %d sec)\n", url.Url, randomNum)
	}

	return url.Url
}

func main() {
	var wg sync.WaitGroup

	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	urls := []WebUrl{
		{"http://www.google.com"},
		{"http://www.Amazon.com"},
		{"http://www.Facebook.com"},
	}

	fmt.Println("Waiting for requests...")

	for _, url := range urls {
		wg.Add(1)
		go func(w WebUrl) {
			defer wg.Done()
			resp := GetResponse(ctx, w)
			fmt.Println(resp)
		}(url)
	}

	wg.Wait() // ✅ call it explicitly, not with defer
}
