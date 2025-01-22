package reporter

import (
	"fmt"
	"sync"
	"time"
)

type Metrics struct {
	mu            sync.Mutex
	StartTime     time.Time
	EndTime       time.Time
	TotalRequests int
	SuccessCount  int
	ErrorCount    int
	TotalDuration time.Duration
	Latency       []time.Duration
}

func (m *Metrics) Begin() {
	m.StartTime = time.Now()
}

func (m *Metrics) End() {
	m.EndTime = time.Now()
}

func (m *Metrics) RecordSuccess(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SuccessCount++
	m.TotalRequests++
	m.TotalDuration += duration
	m.Latency = append(m.Latency, duration)
}

func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorCount++
	m.TotalRequests++
}

func (m *Metrics) Report() {
	m.mu.Lock()
	defer m.mu.Unlock()

	avgLatency := m.TotalDuration / time.Duration(m.TotalRequests)
	totalTime := m.EndTime.Sub(m.StartTime)
	throughput := float64(m.TotalDuration) / totalTime.Seconds()

	fmt.Println("\n--- Load Test Summary ---")
	fmt.Printf("Total Requests: %d\n", m.TotalRequests)
	fmt.Printf("Successful Requests: %d\n", m.SuccessCount)
	fmt.Printf("Failed Requests: %d\n", m.ErrorCount)
	fmt.Printf("Average Latency: %v\n", avgLatency)
	if len(m.Latency) > 0 {
		fmt.Printf("Max Latency: %v\n", maxLatency(m.Latency))
	}
	fmt.Printf("Throughput: %.2f requests/second\n", throughput)
}

func maxLatency(latencies []time.Duration) time.Duration {
	max := latencies[0]
	for _, lat := range latencies {
		if lat > max {
			max = lat
		}
	}
	return max
}
