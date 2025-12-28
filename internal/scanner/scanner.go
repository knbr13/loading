package scanner

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Connection struct {
	LocalAddr  net.IP
	LocalPort  uint16
	RemoteAddr net.IP
	RemotePort uint16
	State      string
	Inode      string
	UID        uint32
}

func ScanTCP(ipv6 bool) ([]Connection, error) {
	path := "/proc/net/tcp"
	if ipv6 {
		path = "/proc/net/tcp6"
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var connections []Connection
	scanner := bufio.NewScanner(file)
	// Skip header
	if scanner.Scan() {
		_ = scanner.Text()
	}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		localAddr, localPort, err := parseAddr(fields[1])
		if err != nil {
			continue
		}
		remoteAddr, remotePort, err := parseAddr(fields[2])
		if err != nil {
			continue
		}

		conn := Connection{
			LocalAddr:  localAddr,
			LocalPort:  localPort,
			RemoteAddr: remoteAddr,
			RemotePort: remotePort,
			State:      fields[3],
			Inode:      fields[9],
		}

		if uid, err := strconv.ParseUint(fields[7], 10, 32); err == nil {
			conn.UID = uint32(uid)
		}

		connections = append(connections, conn)
	}

	return connections, scanner.Err()
}

func parseAddr(s string) (net.IP, uint16, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil, 0, fmt.Errorf("invalid address format")
	}

	ipHex := parts[0]
	portHex := parts[1]

	port, err := strconv.ParseUint(portHex, 16, 16)
	if err != nil {
		return nil, 0, err
	}

	ip, err := decodeHexIP(ipHex)
	if err != nil {
		return nil, 0, err
	}

	return ip, uint16(port), nil
}

func decodeHexIP(h string) (net.IP, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}

	if len(b) == 4 {
		// IPv4: Little Endian
		return net.IP{b[3], b[2], b[1], b[0]}, nil
	} else if len(b) == 16 {
		// IPv6: Each 4-byte block is Little Endian
		ip := make(net.IP, 16)
		for i := 0; i < 4; i++ {
			ip[i*4] = b[i*4+3]
			ip[i*4+1] = b[i*4+2]
			ip[i*4+2] = b[i*4+1]
			ip[i*4+3] = b[i*4]
		}
		return ip, nil
	}

	return nil, fmt.Errorf("invalid IP length: %d", len(b))
}
