package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	recordCount = 500000 // 50万件のテストデータ
)

type DBConfig struct {
	driver         string
	connStr        string
	dbName         string
	createTableSQL string
}

var (
	mysqlConfig = DBConfig{
		driver:  "mysql",
		connStr: "root:password@tcp(localhost:3306)/benchmark_db",
		dbName:  "MySQL",
		createTableSQL: `
			CREATE TABLE IF NOT EXISTS users (
				id_v4 CHAR(36) PRIMARY KEY,
				id_v7 CHAR(36),
				name VARCHAR(255),
				created_at TIMESTAMP
			)
		`,
	}

	postgresConfig = DBConfig{
		driver:  "postgres",
		connStr: "host=localhost port=5432 user=postgres password=password dbname=benchmark_db sslmode=disable",
		dbName:  "PostgreSQL",
		createTableSQL: `
			CREATE TABLE IF NOT EXISTS users (
				id_v4 UUID PRIMARY KEY,
				id_v7 UUID,
				name VARCHAR(255),
				created_at TIMESTAMP
			)
		`,
	}
)

func main() {
	configs := []DBConfig{mysqlConfig, postgresConfig}

	for _, config := range configs {
		db, err := sql.Open(config.driver, config.connStr)
		if err != nil {
			log.Fatalf("Failed to connect to %s: %v", config.dbName, err)
		}
		defer db.Close()

		// テーブル作成
		_, err = db.Exec(config.createTableSQL)
		if err != nil {
			log.Fatalf("Failed to create table in %s: %v", config.dbName, err)
		}

		// データ挿入のベンチマーク
		fmt.Printf("\n=== %s Insert Benchmark ===\n", config.dbName)
		insertBenchmark(db, config)

		// クエリのベンチマーク
		fmt.Printf("\n=== %s Query Benchmark ===\n", config.dbName)
		queryBenchmark(db, config)
	}
}

func insertBenchmark(db *sql.DB, config DBConfig) {
	start := time.Now()

	// バッチサイズを設定
	batchSize := 1000
	valueArgs := make([]interface{}, 0, batchSize*4)

	// プレースホルダーの作成
	var placeholders string
	if config.driver == "postgres" {
		// PostgreSQL用のプレースホルダー ($1, $2, $3, $4), ($5, $6, $7, $8), ...
		placeholder := "($%d, $%d, $%d, $%d)"
		holders := make([]string, batchSize)
		for i := 0; i < batchSize; i++ {
			holders[i] = fmt.Sprintf(placeholder, i*4+1, i*4+2, i*4+3, i*4+4)
		}
		placeholders = strings.Join(holders, ",")
	} else {
		// MySQL用のプレースホルダー (?, ?, ?, ?), (?, ?, ?, ?), ...
		holders := make([]string, batchSize)
		for i := 0; i < batchSize; i++ {
			holders[i] = "(?, ?, ?, ?)"
		}
		placeholders = strings.Join(holders, ",")
	}

	// INSERT文のベース部分
	baseInsertSQL := "INSERT INTO users (id_v4, id_v7, name, created_at) VALUES " + placeholders

	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(baseInsertSQL)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for i := 0; i < recordCount; i++ {
		// UUIDv4の生成
		uuidV4 := uuid.New()

		// UUIDv7の生成
		uuidV7, err := uuid.NewV7()
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to generate UUIDv7: %v", err)
		}

		// バッチ処理用のデータ追加
		valueArgs = append(valueArgs, uuidV4, uuidV7, fmt.Sprintf("user_%d", i), time.Now())

		// バッチサイズに達したらインサート実行
		if len(valueArgs) == batchSize*4 {
			_, err = stmt.Exec(valueArgs...)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Failed to execute batch insert: %v", err)
			}
			valueArgs = valueArgs[:0]
		}
	}

	// 残りのレコードをインサート
	if len(valueArgs) > 0 {
		// 残りのレコード数に合わせてプレースホルダーを調整
		remainingCount := len(valueArgs) / 4
		if config.driver == "postgres" {
			placeholder := "($%d, $%d, $%d, $%d)"
			holders := make([]string, remainingCount)
			for i := 0; i < remainingCount; i++ {
				holders[i] = fmt.Sprintf(placeholder, i*4+1, i*4+2, i*4+3, i*4+4)
			}
			placeholders = strings.Join(holders, ",")
		} else {
			holders := make([]string, remainingCount)
			for i := 0; i < remainingCount; i++ {
				holders[i] = "(?, ?, ?, ?)"
			}
			placeholders = strings.Join(holders, ",")
		}

		finalInsertSQL := "INSERT INTO users (id_v4, id_v7, name, created_at) VALUES " + placeholders
		stmt, err := tx.Prepare(finalInsertSQL)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to prepare final statement: %v", err)
		}

		_, err = stmt.Exec(valueArgs...)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to execute final batch insert: %v", err)
		}
		stmt.Close()
	}

	// トランザクションのコミット
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Inserted %d records in %s\n", recordCount, elapsed)
}

func queryBenchmark(db *sql.DB, config DBConfig) {
	// データベースに応じたプレースホルダーを設定
	placeholder := "?"
	if config.driver == "postgres" {
		placeholder = "$1"
	}

	// UUIDv4でのクエリ実行時間計測
	start := time.Now()
	var count int
	querySQL := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE id_v4 > %s LIMIT 1000", placeholder)
	err := db.QueryRow(querySQL, uuid.Nil).Scan(&count)
	if err != nil {
		log.Printf("Failed to query with UUIDv4: %v", err)
	}
	v4Elapsed := time.Since(start)
	fmt.Printf("UUIDv4 query took: %s\n", v4Elapsed)

	// UUIDv7でのクエリ実行時間計測
	start = time.Now()
	querySQL = fmt.Sprintf("SELECT COUNT(*) FROM users WHERE id_v7 > %s LIMIT 1000", placeholder)
	err = db.QueryRow(querySQL, uuid.Nil).Scan(&count)
	if err != nil {
		log.Printf("Failed to query with UUIDv7: %v", err)
	}
	v7Elapsed := time.Since(start)
	fmt.Printf("UUIDv7 query took: %s\n", v7Elapsed)
}
