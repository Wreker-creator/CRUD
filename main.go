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

	// tasks := []request.Task{
	// 	{ID: 1, Title: "Aman", Description: "First task"},
	// 	{ID: 2, Title: "Akbar", Description: "Second task"},
	// 	{ID: 3, Title: "Anthony", Description: "Third task"},
	// 	{ID: 4, Title: "MKC", Description: "Fourth task"},
	// }

	// store := &request.InMemoryUserStore{
	// 	Tasks: tasks,
	// }
	// server := request.NewTaskServer(store)
	// http.ListenAndServe(":5001", server)

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store := request.NewFileSystemPlayerStore(db)
	server := request.NewTaskServer(store)

	http.ListenAndServe(":5001", server)

}
