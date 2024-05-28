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

func ValidateDB() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS lookup (ip text, port text, status text)")
	if err != nil {
		panic(err)
	}
}

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
		color.Red("FAIL | Test failed")
	}

}

func CreateMemo(args []string) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO lookup (ip, port, status) VALUES (?, ?, ?)", args[0], args[1], args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func UpdateMemo(args []string) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE lookup SET status=? WHERE ip=? AND port=?", args[2], args[0], args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DeletePort(args []string) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM lookup WHERE ip=? AND port=?", args[0], args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DeleteIP(ip string) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM lookup WHERE ip=?", ip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readDb(ip string, memoMap *map[string]string) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := db.Query("SELECT port, status FROM lookup WHERE ip=?", ip)
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
