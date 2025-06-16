package converters

import (
	"fmt"
	"net"
)

func NextIP(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return "", fmt.Errorf("invalid IPv4 address: %s", ipStr)
	}

	// Increment the IP
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break // no overflow, done
		}
	}

	return ip.String(), nil
}
