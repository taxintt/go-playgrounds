package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/taxintt/go-playgrounds/sqlc-dockertest/db"

	"github.com/jackc/pgx/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var q *db.Queries
var conn *pgx.Conn

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	pwd, _ := os.Getwd()

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
		Mounts: []string{
			// docker-entrypoint-initdb.dにschema.sqlをマウントすると、コンテナ起動時に反映される
			fmt.Sprintf("%s/schema.sql:/docker-entrypoint-initdb.d/schema.sql", pwd),
		},
	}, func(config *docker.HostConfig) {
		// 終了時にコンテナを削除する
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	dbPath := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", resource.GetHostPort("5432/tcp"))
	if err := pool.Retry(func() error {
		conn, err = pgx.Connect(context.Background(), dbPath)
		if err != nil {
			return err
		}

		// 接続が確立されているかを確認する
		if conn.Ping(context.Background()); err != nil {
			return err
		}
		q = db.New(conn)

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestUpdateUserAgesWithTransaction(t *testing.T) {
	// ユーザーを作成
	u, err := q.CreateUser(context.Background(), db.CreateUserParams{
		Name:  "test",
		Email: "test@test.com",
		Age:   20,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 年齢を+1する関数を実行
	err = IncrementUserAges(context.Background(), conn, q, u.ID)
	if err != nil {
		t.Fatal(err)
	}

	// 年齢が+1されていることを確認
	q = db.New(conn)
	u, err = q.GetUser(context.Background(), u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if u.Age != 21 {
		t.Fatalf("expected age to be 21, got %d", u.Age)
	}
}
