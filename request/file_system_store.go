package request

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

type FileSystemStore struct {
	database *json.Encoder // we dont need to create a new encoder every time we could just initialise one in our constrcutor

	list List
	/*
		The reason for adding list here is as follows -

		Every time someone calls GetAllTasks() or get specific task, we read the entire file and parse
		it into json. We should not do that BECAUSE FileSystemStore is entirely responsible for the
		state of the league. It should only need to read the file WHEN the program starts, and only update
		when the data changes.

		So we can create a constructor of sorts and store the list as a value.
	*/
}

func NewFileSystemTaskStore(database *os.File) (*FileSystemStore, error) {

	err := initialiseTaskDbFile(database)
	if err != nil {
		return nil, fmt.Errorf("problem initialise the db file %v", err)
	}

	list, err := NewList(database)

	if err != nil {
		return nil, fmt.Errorf("got error %v, while loading tasks from file %s", err, database.Name())
	}

	return &FileSystemStore{
		database: json.NewEncoder(&tape{database}), // encapsulated in tape
		list:     list,
	}, nil

}

func initialiseTaskDbFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("got error %v, while getting info from file %s", err, file.Name())
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	return nil

}

func (f *FileSystemStore) DeleteTask(id int) bool {

	for i := range f.list {
		if f.list[i].ID == id {
			f.list = append(f.list[:i], f.list[i+1:]...)
			// f.database.Seek(0, io.SeekStart)
			// json.NewEncoder(f.database).Encode(&f.list)
			f.database.Encode(&f.list)
			return true
		}
	}
	return false
}

func (f *FileSystemStore) GetAllTasks() List {

	sort.Slice(f.list, func(i, j int) bool {
		return f.list[i].ID < f.list[j].ID
	})

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
			// json.NewEncoder(f.database).Encode(&f.list)
			f.database.Encode(&f.list)
			return true
		}
	}

	return false

}

func (f *FileSystemStore) AddTask(task Task) {
	f.list = append(f.list, task)
	// f.database.Seek(0, io.SeekStart)
	// json.NewEncoder(f.database).Encode(&f.list)
	f.database.Encode(&f.list)
}
