package state

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/knbr13/glow-net/internal/enrichment"
	"github.com/knbr13/glow-net/internal/mapper"
	"github.com/knbr13/glow-net/internal/scanner"
)

type RowData struct {
	PID          int
	ProcessName  string
	LocalAddr    string
	RemoteAddr   string
	RemoteHost   string
	Country      string
	DownloadRate float64
	UploadRate   float64
}

type AppState struct {
	sync.RWMutex
	Connections []scanner.Connection
	ProcessMap  map[string]mapper.ProcessInfo
	GlobalStats scanner.GlobalStats
	Enricher    *enrichment.Enricher
	LastUpdate  time.Time

	// Internal stats for rate calculation
	prevProcInodes map[string]bool // placeholder for now
}

func NewAppState(enricher *enrichment.Enricher) *AppState {
	return &AppState{
		Enricher:   enricher,
		ProcessMap: make(map[string]mapper.ProcessInfo),
	}
}

func (s *AppState) Refresh(bwScanner *scanner.BandwidthScanner) error {
	s.Lock()
	defer s.Unlock()

	// 1. Scan connections
	conns, err := scanner.ScanTCP(false) // IPv4
	if err == nil {
		s.Connections = conns
	}
	conns6, err := scanner.ScanTCP(true) // IPv6
	if err == nil {
		s.Connections = append(s.Connections, conns6...)
	}

	// 2. Map processes
	pm, err := mapper.GetInodeToProcessMap()
	if err == nil {
		s.ProcessMap = pm
	}

	// 3. Global stats
	gs, err := bwScanner.GetGlobalStats()
	if err == nil {
		s.GlobalStats = gs
	}

	s.LastUpdate = time.Now()
	return nil
}

func (s *AppState) GetTableRows() []RowData {
	s.RLock()
	defer s.RUnlock()

	var rows []RowData
	for _, conn := range s.Connections {
		proc, ok := s.ProcessMap[conn.Inode]
		pid := -1
		name := "unknown"
		if ok {
			pid = proc.PID
			name = proc.Name
		}

		country := s.Enricher.LookupCountry(conn.RemoteAddr)
		remoteHost := s.Enricher.LookupDNS(conn.RemoteAddr)

		rows = append(rows, RowData{
			PID:         pid,
			ProcessName: name,
			LocalAddr:   net.JoinHostPort(conn.LocalAddr.String(), fmt.Sprintf("%d", conn.LocalPort)),
			RemoteAddr:  net.JoinHostPort(conn.RemoteAddr.String(), fmt.Sprintf("%d", conn.RemotePort)),
			RemoteHost:  remoteHost,
			Country:     country,
		})
	}
	return rows
}
