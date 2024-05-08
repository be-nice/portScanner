package memo

import (
	"database/sql"
	"fmt"
	"os"
	"portscan/scanner"
	util "portscan/utility"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
)

func GetMemo(ip string) {
	memoMap := make(map[string]string)
	readDb(ip, &memoMap)
	var status string
	tesStatus := true
	for key, val := range memoMap {
		scan := util.Scan{
			IP:        ip,
			StartPort: key,
		}
		err := scanner.ScanPorts(scan)
		if err != nil {
			status = "closed"
		} else {
			status = "open"
		}
		if (err == nil && val == "open") || (err != nil && val == "closed") {
			str := color.GreenString(fmt.Sprintf("Port %s expected %s | status: %s | SUCCESS", key, val, status))
			fmt.Println(str)
		} else {
			tesStatus = false
			str := color.RedString(fmt.Sprintf("Port %s expected %s | status: %s | FAIL", key, val, status))
			fmt.Println(str)
		}
	}

	if tesStatus {
		color.Green("SUCCESS, All tests passed")
	} else {
		color.Red("Fail | All tests did not pass")
	}

}

func readDb(ip string, memoMap *map[string]string) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := db.Query("select port, status from lookup where ip=?", ip)
	db.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var port string
	var val string
	for rows.Next() {
		rows.Scan(&port, &val)
		(*memoMap)[port] = val
	}
	defer rows.Close()
}
