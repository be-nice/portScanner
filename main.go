package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"portscan/scanner"
	util "portscan/utility"
)

func main() {
	wg := sync.WaitGroup{}

	scan, err := validateArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	sPort, err := strconv.Atoi(scan.StartPort)
	if err != nil {
		panic(err)
	}
	ePort, err := strconv.Atoi(scan.EndPort)
	if err != nil {
		panic(err)
	}
	for i := sPort; i <= ePort; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			scan.StartPort = strconv.Itoa(port)
			scanner.ScanPorts(scan)
		}(i)
	}
	wg.Wait()
}

func validateArgs() (util.Scan, error) {
	if len(os.Args) < 2 {
		return util.Scan{}, errors.New("no arguments provided")
	}

	ip := os.Args[1]
	if !util.ValidIP(ip) {
		return util.Scan{}, errors.New("invalid IP address")
	}

	if len(os.Args) == 2 {
		return util.Scan{IP: ip, StartPort: "1", EndPort: "65535"}, nil
	}

	if len(os.Args) != 3 && len(os.Args) != 4 {
		return util.Scan{}, errors.New("invalid number of arguments")
	}

	startPort := os.Args[2]
	if !util.ValidPort(startPort) {
		return util.Scan{}, errors.New("invalid start port")
	}

	if len(os.Args) == 3 {
		return util.Scan{IP: ip, StartPort: startPort, EndPort: startPort}, nil
	}

	endPort := os.Args[3]
	if !util.ValidPort(endPort) {
		return util.Scan{}, errors.New("invalid end port")
	}

	return util.Scan{IP: ip, StartPort: startPort, EndPort: endPort}, nil
}
