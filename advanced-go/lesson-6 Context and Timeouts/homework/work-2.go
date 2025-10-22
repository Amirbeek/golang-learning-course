package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Website struct {
	url string
}

func CheckWebsite(ctx context.Context, w Website) error {

	time.Sleep(1 * time.Second)
	select {
	case <-ctx.Done():
		fmt.Println(w.url, "check canceled:", ctx.Err())
		return ctx.Err()
	case <-time.After(time.Second * 3):
		fmt.Println(w.url, "is healthy ✅")
	}
	return nil
}

func main() {
	var wg sync.WaitGroup
	defer wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	webs := []Website{
		{url: "http://www.google.com"},
		{url: "http://www.Facebook.com"},
		{url: "http://www.GitHub.com"},
	}
	for _, web := range webs {
		wg.Add(1)
		go func(w Website) {
			defer wg.Done()
			if err := CheckWebsite(ctx, w); err != nil {
				fmt.Println("Website check failed:", err)
			}
		}(web)

	}

	wg.Wait()
	print("All Websites are healthy")
}
