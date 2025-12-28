package scanner

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type GlobalStats struct {
	DownloadSpeed float64 // bytes/sec
	UploadSpeed   float64 // bytes/sec
}

type BandwidthScanner struct {
	prevRx uint64
	prevTx uint64
}

func (s *BandwidthScanner) GetGlobalStats() (GlobalStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return GlobalStats{}, err
	}
	defer file.Close()

	var totalRx, totalTx uint64
	scanner := bufio.NewScanner(file)
	// Skip 2 header lines
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 10 {
			continue
		}
		// fields[1] is Rx bytes, fields[9] is Tx bytes
		rx, _ := strconv.ParseUint(fields[1], 10, 64)
		tx, _ := strconv.ParseUint(fields[9], 10, 64)
		totalRx += rx
		totalTx += tx
	}

	var stats GlobalStats
	if s.prevRx > 0 {
		stats.DownloadSpeed = float64(totalRx - s.prevRx)
		stats.UploadSpeed = float64(totalTx - s.prevTx)
	}

	s.prevRx = totalRx
	s.prevTx = totalTx

	return stats, nil
}
