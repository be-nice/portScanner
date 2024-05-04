package scanner

import (
	"net"
	"time"

	util "portscan/utility"
)

func ScanPorts(scan util.Scan) error {
	conn, err := net.DialTimeout("tcp", scan.IP+":"+scan.StartPort, 3*time.Second)
	if err != nil {
		return err
	} else {
		conn.Close()
		return nil
	}
}
