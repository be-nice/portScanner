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

// DB interface for database operations
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Close() error
}

// DBClient struct to hold the database connection
type DBClient struct {
	db *sql.DB
}

// NewDBClient initializes and returns a new DBClient
func NewDBClient(dataSourceName string) (*DBClient, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DBClient{db: db}, nil
}

// Exec executes a query without returning any rows
func (client *DBClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	return client.db.Exec(query, args...)
}

// Query executes a query that returns rows
func (client *DBClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return client.db.Query(query, args...)
}

// Close closes the database connection
func (client *DBClient) Close() error {
	return client.db.Close()
}

// ValidateDB ensures the lookup table exists
func ValidateDB(client DB) {
	_, err := client.Exec("CREATE TABLE IF NOT EXISTS lookup (ip text, port text, status text)")
	if err != nil {
		panic(err)
	}
}

// GetMemo retrieves and scans ports from the database
func GetMemo(client DB, ip string) {
	memoMap := make(map[string]string)
	readDb(client, ip, &memoMap)
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

// CreateMemo inserts a new record into the database
func CreateMemo(client DB, args []string) {
	_, err := client.Exec("INSERT INTO lookup (ip, port, status) VALUES (?, ?, ?)", args[0], args[1], args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// UpdateMemo updates an existing record in the database
func UpdateMemo(client DB, args []string) {
	_, err := client.Exec("UPDATE lookup SET status=? WHERE ip=? AND port=?", args[2], args[0], args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// DeletePort deletes a port record from the database
func DeletePort(client DB, args []string) {
	_, err := client.Exec("DELETE FROM lookup WHERE ip=? AND port=?", args[0], args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// DeleteIP deletes all records for a given IP from the database
func DeleteIP(client DB, ip string) {
	_, err := client.Exec("DELETE FROM lookup WHERE ip=?", ip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// readDb reads records from the database into a map
func readDb(client DB, ip string, memoMap *map[string]string) {
	rows, err := client.Query("SELECT port, status FROM lookup WHERE ip=?", ip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rows.Close()

	var port string
	var val string
	for rows.Next() {
		rows.Scan(&port, &val)
		(*memoMap)[port] = val
	}
}
