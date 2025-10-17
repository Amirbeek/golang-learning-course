package main

import (
	"fmt"
	"sync"
	"time"
)

func downloadFile(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Downloading file:", name)
	time.Sleep(time.Second * 2)
	fmt.Println("Finished downloading file:", name)
}

func main() {
	files := []string{"file1.zip", "file2.mp4", "file3.pdf", "file4.jpg"}
	start := time.Now()

	var wg sync.WaitGroup
	for _, f := range files {
		wg.Add(1)
		go downloadFile(f, &wg)
	}
	wg.Wait()
	fmt.Println("Finished Downloading file on ", time.Since(start))
	fmt.Println("Finished downloading files:", files)
}
