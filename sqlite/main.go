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

	days := 2
	number := 10

	// sqlStmt := fmt.Sprintf(`
	// SELECT datetime(last_visit_time/1000000-11644473600, "unixepoch", "localtime") as last_visited, url, title
	// FROM urls
	// WHERE last_visited > datetime('now', '-%d days')
	// ORDER BY last_visited DESC
	// LIMIT %d;
	// `, days, number)

	sqlStmt := fmt.Sprintf(`
	SELECT datetime(last_visit_time/1000000-11644473600, "unixepoch", "localtime") as last_visited, last_visit_time, url, title 
	FROM urls
	WHERE last_visited > datetime('now', '-%d days')
	ORDER BY last_visited DESC
	LIMIT %d;
	`, days, number)

	fmt.Println(sqlStmt)
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var last_visited, title, url string
		var last_visit_time int64
		err = rows.Scan(&last_visited, &last_visit_time, &title, &url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(last_visit_time, title, " | ", url)
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
