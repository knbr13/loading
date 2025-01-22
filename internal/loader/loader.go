package loader

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/knbr13/loading/internal/reporter"
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
	metrics := &reporter.Metrics{}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	metrics.Begin()

	for i := 0; i < options.RequestCount; i++ {
		wg.Add(1)
		ch <- 1

		go func(requestID int) {
			defer wg.Done()
			start := time.Now()

			req, err := http.NewRequest(options.Method, options.URL, nil)
			if err != nil {
				fmt.Printf("Request %d failed to create: %v\n", requestID, err)
				metrics.RecordError()
				<-ch
				return
			}

			for key, value := range options.Headers {
				req.Header.Set(key, value)
			}

			resp, err := client.Do(req)
			duration := time.Since(start)
			if err != nil {
				fmt.Printf("Request %d failed: %v\n", requestID, err)
				metrics.RecordError()
			} else {
				metrics.RecordSuccess(duration)
				resp.Body.Close()
			}

			<-ch
		}(i + 1)
	}

	wg.Wait()
	metrics.End()
	metrics.Report()
}
