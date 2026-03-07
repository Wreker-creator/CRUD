package request

import (
	"io"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("List from reader", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"ID": 1, "Title": "A", "Description": "First Task"},
			{"ID": 2, "Title": "B", "Description": "Second Task"},
			{"ID": 3, "Title": "C", "Description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{Database: database}
		store := NewFileSystemPlayerStore(database)
		got := store.GetAllTasks()

		want := []Task{
			{ID: 1, Title: "A", Description: "First Task"},
			{ID: 2, Title: "B", Description: "Second Task"},
			{ID: 3, Title: "C", Description: "Third Task"},
		}

		assertTasks(t, got, want)

		got = store.GetAllTasks()
		assertTasks(t, got, want)

	})

	t.Run("get a specific task from an id", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"ID": 1, "Title": "A", "Description": "First Task"},
			{"ID": 2, "Title": "B", "Description": "Second Task"},
			{"ID": 3, "Title": "C", "Description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{database}
		store := NewFileSystemPlayerStore(database)
		got, _ := store.GetAllTasks().Find(1)

		want := Task{ID: 1, Title: "A", Description: "First Task"}
		assertTask(t, *got, want)

	})

	t.Run("update existing task", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"ID": 1, "Title": "A", "Description": "First Task"},
			{"ID": 2, "Title": "B", "Description": "Second Task"},
			{"ID": 3, "Title": "C", "Description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{Database: database}
		store := NewFileSystemPlayerStore(database)

		task := Task{ID: 2, Title: "A", Description: "First Task"}
		store.UpdateTask(2, task)

		got, _ := store.GetAllTasks().Find(2)

		assertTask(t, *got, task)

	})

	t.Run("store new tasks", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"ID": 1, "Title": "A", "Description": "First Task"},
			{"ID": 2, "Title": "B", "Description": "Second Task"},
			{"ID": 3, "Title": "C", "Description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{Database: database}
		store := NewFileSystemPlayerStore(database)
		task := Task{ID: 4, Title: "D", Description: "Fourth Task"}
		store.AddTask(task)

		got, _ := store.GetAllTasks().Find(4)

		assertTask(t, *got, task)

	})

	t.Run("Delete task", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"ID": 1, "Title": "A", "Description": "First Task"},
			{"ID": 2, "Title": "B", "Description": "Second Task"},
			{"ID": 3, "Title": "C", "Description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{Database: database}
		store := NewFileSystemPlayerStore(database)

		id := 3
		store.DeleteTask(id)

		task, _ := store.GetAllTasks().Find(id)
		assertTask(t, *task, Task{})

		got := store.GetAllTasks()
		want := []Task{
			{ID: 1, Title: "A", Description: "First Task"},
			{ID: 2, Title: "B", Description: "Second Task"},
		}

		assertTasks(t, got, want)

	})

}

func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("unable to create a temp file createTempFile func, %v", err)
	}

	tempFile.Write([]byte(initialData))

	removeFile := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	return tempFile, removeFile

}
