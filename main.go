package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Scan struct {
	ip        string
	startPort string
	endPort   string
}

var commmonPorts = map[string]string{
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

func main() {
	wg := sync.WaitGroup{}

	scan, err := validateArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	sPort, err := strconv.Atoi(scan.startPort)
	if err != nil {
		panic(err)
	}
	ePort, err := strconv.Atoi(scan.endPort)
	if err != nil {
		panic(err)
	}

	for i := sPort; i <= ePort; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scan.startPort = strconv.Itoa(i)
			scanPorts(scan)
		}()
	}
	wg.Wait()
}

func validateArgs() (Scan, error) {
	if len(os.Args) < 2 {
		return Scan{}, errors.New("no arguments provided")
	}

	ip := os.Args[1]
	if !validIP(ip) {
		return Scan{}, errors.New("invalid IP address")
	}

	if len(os.Args) == 2 {
		return Scan{ip, "1", "65535"}, nil
	}

	if len(os.Args) != 3 && len(os.Args) != 4 {
		return Scan{}, errors.New("invalid number of arguments")
	}

	startPort := os.Args[2]
	if !validPort(startPort) {
		return Scan{}, errors.New("invalid start port")
	}

	if len(os.Args) == 3 {
		return Scan{ip, startPort, startPort}, nil
	}

	endPort := os.Args[3]
	if !validPort(endPort) {
		return Scan{}, errors.New("invalid end port")
	}

	return Scan{ip, startPort, endPort}, nil
}

func scanPorts(scan Scan) {
	conn, err := net.DialTimeout("tcp", scan.ip+":"+scan.startPort, 2*time.Second)
	if err != nil {
		return
	} else {
		if val, ok := commmonPorts[scan.startPort]; ok {
			str := fmt.Sprintf("Port %s is open | (%s)", scan.startPort, val)
			color.Green(str)
		} else {
			fmt.Println("Port", scan.startPort, "is open")
		}
		conn.Close()
	}
}

func validIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func validPort(port string) bool {
	_, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	if port <= "0" || port > "65535" {
		return false
	}
	return true
}
