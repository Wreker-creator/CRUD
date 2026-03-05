package main

import (
	"bookstore/request"
	"net/http"
)

func main() {

	tasks := []request.Task{
		{ID: 1, Title: "Aman", Description: "First task"},
		{ID: 2, Title: "Akbar", Description: "Second task"},
		{ID: 3, Title: "Anthony", Description: "Third task"},
		{ID: 4, Title: "MKC", Description: "Fourth task"},
	}

	store := &request.InMemoryUserStore{
		Tasks: tasks,
	}
	server := request.NewTaskServer(store)
	http.ListenAndServe(":5001", server)

}
