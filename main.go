package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"portscan/memo"
	"portscan/scanner"
	util "portscan/utility"

	"github.com/fatih/color"
)

var portMap = make(map[string]string)
var wg sync.WaitGroup

func main() {
	scan, err := validateArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println("made 3")
	file, err := os.Open("portlist.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanb := bufio.NewScanner(file)
	for scanb.Scan() {
		str := strings.Fields(scanb.Text())
		portMap[strings.Split(str[1], "/")[0]] = str[0]
	}
	if scan.DefaultScan {
		for key, val := range portMap {
			scan.StartPort = key
			scan.Service = val
			wg.Add(1)
			go write(scan)
		}
	} else {
		sPort, err := strconv.Atoi(scan.StartPort)
		if err != nil {
			panic(err)
		}
		ePort, err := strconv.Atoi(scan.EndPort)
		if err != nil {
			panic(err)
		}
		for i := sPort; i <= ePort; i++ {
			port := strconv.Itoa(i)
			if val, ok := portMap[port]; ok {
				scan.Service = val
			}
			scan.StartPort = port
			wg.Add(1)
			go write(scan)
		}
	}
	wg.Wait()
}

func write(scan util.Scan) {
	defer wg.Done()
	err := scanner.ScanPorts(scan)
	if err != nil {
		return
	}
	if scan.Service != "" {
		str := fmt.Sprintf("Port %s is open | (%s)", scan.StartPort, scan.Service)
		color.Green(str)
	} else {
		fmt.Println("Port", scan.StartPort, "is open")
	}
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
		return util.Scan{IP: ip, DefaultScan: true, StartPort: "1", EndPort: "65535"}, nil
	}

	if len(os.Args) != 3 && len(os.Args) != 4 {
		return util.Scan{}, errors.New("invalid number of arguments")
	}
	if len(os.Args) == 3 && os.Args[2] == "-t" {
		memo.GetMemo(os.Args[1])
		os.Exit(0)
	}

	startPort := os.Args[2]
	if !util.ValidPort(startPort) {
		return util.Scan{}, errors.New("invalid start port")
	}

	if len(os.Args) == 3 {
		return util.Scan{IP: ip, DefaultScan: false, StartPort: startPort, EndPort: startPort}, nil
	}

	endPort := os.Args[3]
	if !util.ValidPort(endPort) {
		return util.Scan{}, errors.New("invalid end port")
	}

	return util.Scan{IP: ip, DefaultScan: false, StartPort: startPort, EndPort: endPort}, nil
}
