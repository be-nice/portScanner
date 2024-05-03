package scanner

import (
	"fmt"
	"net"
	"time"

	util "portscan/utility"

	"github.com/fatih/color"
)

func ScanPorts(scan util.Scan) {
	conn, err := net.DialTimeout("tcp", scan.IP+":"+scan.StartPort, 3*time.Second)
	if err != nil {
		return
	} else {
		if scan.Service != "" {
			str := fmt.Sprintf("Port %s is open | (%s)", scan.StartPort, scan.Service)
			color.Green(str)
		} else {
			fmt.Println("Port", scan.StartPort, "is open")
		}
		conn.Close()
	}
}
