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

func init() {
	memo.ValidateDB()
}

func main() {
	scan, err := validateArgs()
	if err != nil {
		fmt.Println(err)
		fmt.Println("type 'help' to see usage")
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
	if len(os.Args) < 2 || len(os.Args) > 5 {
		return util.Scan{}, errors.New("invalid number of arguments")
	}

	args := os.Args[1:]
	var ip string

	if strings.Split(args[0], "")[0] == "-" {
		ip = args[1]
	} else {
		ip = args[0]
	}

	if !util.ValidIP(ip) {
		return util.Scan{}, errors.New("invalid IP address")
	}

	if args[0] != ip {
		if len(args) != 2 && len(args) != 4 {
			return util.Scan{}, errors.New("invalid number of arguments")
		}
		if len(args) == 4 {
			if !util.ValidPort(args[2]) {
				return util.Scan{}, errors.New("invalid port")
			}
			if args[3] != "open" && args[3] != "closed" {
				return util.Scan{}, errors.New("invalid operation")
			}
		}
		switch args[0] {
		case "-t": // run tests
			memo.GetMemo(args[1])
			os.Exit(0)
		case "-c": // create db entry
			memo.CreateMemo(args[1:])
			os.Exit(0)
		case "-e": // edit db entry
			memo.UpdateMemo(args[1:])
			os.Exit(0)
		default:
			return util.Scan{}, errors.New("invalid flag")
		}
	} else {
		switch {
		case len(args) == 1:
			return util.Scan{IP: ip, DefaultScan: true, StartPort: "1", EndPort: "65535"}, nil
		case len(args) == 2:
			if !util.ValidPort(args[1]) {
				return util.Scan{}, errors.New("invalid port")
			}
			return util.Scan{IP: ip, DefaultScan: false, StartPort: args[1], EndPort: args[1]}, nil
		case len(args) == 3:
			if !util.ValidPort(args[1]) || !util.ValidPort(args[2]) {
				return util.Scan{}, errors.New("invlid port")
			}
			return util.Scan{IP: ip, DefaultScan: false, StartPort: args[1], EndPort: args[2]}, nil
		default:
			return util.Scan{}, errors.New("invalid number of arguments")
		}
	}
	return util.Scan{}, errors.New("unkown error")
}
