package reporter

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
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

	totalTime := m.EndTime.Sub(m.StartTime)
	throughput := float64(m.TotalRequests) / totalTime.Seconds()
	avgLatency := m.TotalDuration / time.Duration(m.TotalRequests)

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	fmt.Println("\n--- Load Test Summary ---")
	fmt.Printf("Total Requests: %s\n", blue(m.TotalRequests))
	fmt.Printf("Successful Requests: %s\n", green(m.SuccessCount))
	fmt.Printf("Failed Requests: %s\n", red(m.ErrorCount))
	fmt.Printf("Average Latency: %s\n", blue(avgLatency))
	if len(m.Latency) > 0 {
		fmt.Printf("Max Latency: %s\n", blue(maxLatency(m.Latency)))
	}
	fmt.Printf("Throughput: %s requests/second\n", blue(fmt.Sprintf("%.2f", throughput)))
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
