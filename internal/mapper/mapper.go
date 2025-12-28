package mapper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProcessInfo struct {
	PID  int
	Name string
}

func GetInodeToProcessMap() (map[string]ProcessInfo, error) {
	result := make(map[string]ProcessInfo)

	dirs, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}

		pid := 0
		_, err := fmt.Sscanf(d.Name(), "%d", &pid)
		if err != nil {
			continue
		}

		name := getProcessName(pid)
		inodes := getProcessInodes(pid)

		for _, inode := range inodes {
			result[inode] = ProcessInfo{
				PID:  pid,
				Name: name,
			}
		}
	}

	return result, nil
}

func getProcessName(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func getProcessInodes(pid int) []string {
	var inodes []string
	fdPath := fmt.Sprintf("/proc/%d/fd", pid)
	fds, err := os.ReadDir(fdPath)
	if err != nil {
		return nil
	}

	for _, fd := range fds {
		link, err := os.Readlink(filepath.Join(fdPath, fd.Name()))
		if err != nil {
			continue
		}

		if strings.HasPrefix(link, "socket:[") && strings.HasSuffix(link, "]") {
			inode := link[8 : len(link)-1]
			inodes = append(inodes, inode)
		}
	}

	return inodes
}
