package util

import (
	"net"
	"strconv"
)

var CommmonPorts = map[string]string{
	"20":   "FTP",
	"21":   "FTP control",
	"22":   "SSH",
	"23":   "Telnet",
	"25":   "SMTP",
	"53":   "DNS",
	"80":   "HTTP",
	"110":  "POP3",
	"143":  "IMAP",
	"443":  "HTTPS",
	"465":  "SMTPS",
	"591":  "HTTP alternate",
	"993":  "IMAPS",
	"995":  "POP3S",
	"3306": "MySQL",
	"5432": "PostgreSQL",
	"8008": "HTTP alternate",
	"8080": "HTTP alternate",
	"8443": "HTTPS alternate",
}

type Scan struct {
	IP        string
	StartPort string
	EndPort   string
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
