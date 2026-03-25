package main

import (
	"log"
	"net/http"
	"os"
	"task/task"
)

func main() {
	// os.Getenv reads the DATABASE_URL variable injected by docker-compose.
	// If it's empty, we fail immediately with a clear message rather than
	// crashing later with a confusing nil pointer or connection error.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	store, err := task.NewPostgresTaskStore(dsn)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	server := task.NewTaskServer(store)

	log.Println("Task Manager API running on :5001")

	if err := http.ListenAndServe(":5001", server); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
