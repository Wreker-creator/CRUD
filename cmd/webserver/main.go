package main

import (
	"log"
	"net/http"
	"os"
	"task/task"
)

// "bookstore/request"
// "log"
// "os"

const dbFileName = "task.db.json"

func main() {

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := task.NewFileSystemTaskStore(db)
	if err != nil {
		log.Fatalf("got error %v", err)
	}
	server := task.NewTaskServer(store)

	http.ListenAndServe(":5001", server)

}
