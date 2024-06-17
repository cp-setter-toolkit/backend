// nolint
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/thepluck/cp-setter-toolkit/ent"
)

// getClient opens a connection to the database.
func getClient() (*ent.Client, error) {
	return ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
	))
}

func main() {
	var client *ent.Client
	var connected bool
	// Wait and retry until the database is ready.
	for start := time.Now(); time.Since(start) < 30 * time.Second; {
		client, err := getClient()
		if err == nil {
			connected = true
			break
		}
	}
	if !connected {
		log.Fatalf("failed to connect to the database")
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
