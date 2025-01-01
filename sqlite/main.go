package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./history.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	SELECT datetime(last_visit_time/1000000-11644473600, "unixepoch", "localtime") as last_visited, url, title 
	FROM urls
	order by last_visited desc 
	limit 10;
	`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var last_visited, title, url string
		err = rows.Scan(&last_visited, &title, &url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(last_visited, title, " | ", url)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

// func fetchHistoryFileLocation() (string, error) {
// 	os := runtime.GOOS
// 	switch os {
// 	case "windows":
// 		return `AppData\Local\Google\Chrome\User Data\Default\history`, nil
// 	case "darwin":
// 		return `~/Library/Application\ Support/Google/Chrome/Default/History`, nil
// 	default:
// 		return "", fmt.Errorf("OS not supported")
// 	}
// }
