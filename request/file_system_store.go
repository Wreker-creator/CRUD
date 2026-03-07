package request

import (
	"encoding/json"
	"io"
	"os"
)

type FileSystemStore struct {
	database io.Writer // we are going to encapsulate in tape.
	list     List
	/*
		The reason for adding list here is as follows -

		Every time someone calls GetAllTasks() or get specific task, we read the entire file and parse
		it into json. We should not do that BECAUSE FileSystemStore is entirely responsible for the
		state of the league. It should only need to read the file WHEN the program starts, and only update
		when the data changes.

		So we can create a constructor of sorts and store the list as a value.
	*/
}

func NewFileSystemTaskStore(database *os.File) *FileSystemStore {

	database.Seek(0, io.SeekStart)
	list, _ := NewList(database)

	return &FileSystemStore{
		database: &tape{database}, // encapsulated in tape
		list:     list,
	}

}

func (f *FileSystemStore) DeleteTask(id int) bool {

	for i := range f.list {
		if f.list[i].ID == id {
			f.list = append(f.list[:i], f.list[i+1:]...)
			// f.database.Seek(0, io.SeekStart)
			json.NewEncoder(f.database).Encode(&f.list)
			return true
		}
	}
	return false
}

func (f *FileSystemStore) GetAllTasks() List {
	return f.list
}

// func (f *FileSystemStore) GetSepcificTask(id int) (Task, bool) {

// 	task := f.GetAllTasks().Find(id)
// 	if task != nil {
// 		return *task, true
// 	}
// 	return Task{}, false

// }

func (f *FileSystemStore) UpdateTask(id int, task Task) bool {

	for i := range f.list {
		if f.list[i].ID == id {
			f.list[i] = task
			// f.database.Seek(0, io.SeekStart)
			json.NewEncoder(f.database).Encode(&f.list)
			return true
		}
	}

	return false

}

func (f *FileSystemStore) AddTask(task Task) {
	f.list = append(f.list, task)
	// f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(&f.list)
}
