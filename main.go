package main

import (
	"bookstore/request"
	"log"
	"net/http"
	"os"
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

	store := request.NewFileSystemTaskStore(db)
	server := request.NewTaskServer(store)

	http.ListenAndServe(":5001", server)

}
