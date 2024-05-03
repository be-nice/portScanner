package util

import (
	"net"
	"strconv"
)

type Scan struct {
	IP          string
	DefaultScan bool
	StartPort   string
	EndPort     string
	Service     string
}

func ValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func ValidPort(port string) bool {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	if intPort <= 0 || intPort > 65535 {
		return false
	}
	return true
}
