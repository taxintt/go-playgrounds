package main

import (
	"context"
	"log"

	"github.com/taxintt/go-playgrounds/ent-with-atlas/ent"

	_ "github.com/lib/pq"
)

func main() {
	client, err := ent.Open("postgres", "host=localhost port=5432 user=your_username dbname=your_database password=your_password")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
