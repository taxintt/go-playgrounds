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
	select substr(title, 1, $cols) as title, url
  from urls order by last_visit_time desc limit 10;
	`

	rows, err := db.Query(sqlStmt, 50)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var title, url string
		err = rows.Scan(&title, &url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(title, " | ", url)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
