package loader

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type RequestOptions struct {
	Method       string
	URL          string
	Headers      map[string]string
	Concurrency  int
	RequestCount int
}

func LoadTest(options RequestOptions) {
	var wg sync.WaitGroup
	ch := make(chan int, options.Concurrency)

	for i := 0; i < options.RequestCount; i++ {
		wg.Add(1)
		ch <- 1

		go func(requestID int) {
			defer wg.Done()
			start := time.Now()

			req, err := http.NewRequest(options.Method, options.URL, nil)
			if err != nil {
				fmt.Printf("Request %d failed to create: %v\n", requestID, err)
				<-ch
				return
			}

			for key, value := range options.Headers {
				req.Header.Set(key, value)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("Request %d failed: %v\n", requestID, err)
			} else {
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Request %d completed in %v. Status: %s\n", requestID, time.Since(start), resp.Status)
				resp.Body.Close()
				_ = body
			}

			<-ch
		}(i + 1)
	}

	wg.Wait()
	fmt.Println("load test completed!")
}
