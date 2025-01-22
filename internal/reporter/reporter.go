package reporter

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Value"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("*")
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	data := [][]string{
		{"Total Requests", strconv.Itoa(m.TotalRequests)},
		{"Successful Requests", color.New(color.FgGreen).Sprint(m.SuccessCount)},
		{"Failed Requests", color.New(color.FgRed).Sprint(m.ErrorCount)},
		{"Average Latency", avgLatency.String()},
		{"Max Latency", maxLatency(m.Latency).String()},
		{"Throughput (req/sec)", strconv.FormatFloat(throughput, 'f', 2, 64)},
	}

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Println(cyan("\nTest completed successfully! Analyze the results above.\n"))
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
