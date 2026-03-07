package request

import (
	"io"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("List from reader", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"id": 1, "title": "A", "description": "First Task"},
			{"id": 2, "title": "B", "description": "Second Task"},
			{"id": 3, "title": "C", "description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{Database: database}
		store, err := NewFileSystemTaskStore(database)
		assertNoError(t, err)

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
			{"id": 1, "title": "A", "description": "First Task"},
			{"id": 2, "title": "B", "description": "Second Task"},
			{"id": 3, "title": "C", "description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{database}
		store, err := NewFileSystemTaskStore(database)
		assertNoError(t, err)
		got, _ := store.GetAllTasks().Find(1)

		want := Task{ID: 1, Title: "A", Description: "First Task"}
		assertTask(t, *got, want)

	})

	t.Run("update existing task", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"id": 1, "title": "A", "description": "First Task"},
			{"id": 2, "title": "B", "description": "Second Task"},
			{"id": 3, "title": "C", "description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{Database: database}
		store, err := NewFileSystemTaskStore(database)
		assertNoError(t, err)

		task := Task{ID: 2, Title: "A", Description: "First Task"}
		store.UpdateTask(2, task)

		got, _ := store.GetAllTasks().Find(2)

		assertTask(t, *got, task)

	})

	t.Run("store new tasks", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"id": 1, "title": "A", "description": "First Task"},
			{"id": 2, "title": "B", "description": "Second Task"},
			{"id": 3, "title": "C", "description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{Database: database}
		store, err := NewFileSystemTaskStore(database)
		assertNoError(t, err)
		task := Task{ID: 4, Title: "D", Description: "Fourth Task"}
		store.AddTask(task)

		got, _ := store.GetAllTasks().Find(4)

		assertTask(t, *got, task)

	})

	t.Run("Delete task", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"id": 1, "title": "A", "description": "First Task"},
			{"id": 2, "title": "B", "description": "Second Task"},
			{"id": 3, "title": "C", "description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{Database: database}
		store, err := NewFileSystemTaskStore(database)
		assertNoError(t, err)

		id := 3
		store.DeleteTask(id)

		task, task_err := store.GetAllTasks().Find(id)
		if task != nil {
			t.Errorf("expected task to be deleted, but found %v", task_err)
		}

		got := store.GetAllTasks()
		want := []Task{
			{ID: 1, Title: "A", Description: "First Task"},
			{ID: 2, Title: "B", Description: "Second Task"},
		}

		assertTasks(t, got, want)

	})

	t.Run("tasks are sorted by id", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"id": 2, "title": "B", "description": "Second Task"},
			{"id": 1, "title": "A", "description": "First Task"},
			{"id": 3, "title": "C", "description": "Third Task"}
		]`)

		defer cleanDatabase()

		// store := FileSystemStore{Database: database}
		store, err := NewFileSystemTaskStore(database)
		assertNoError(t, err)

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

}

// changed to os.File because we are wrapping the database into tape
// which allows us to truncate.
func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("unable to create a temp file createTempFile func, %v", err)
	}

	tempFile.Write([]byte(initialData))
	tempFile.Seek(0, io.SeekStart)
	// bringing the cursor to the start
	// once the data has been written.

	removeFile := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	return tempFile, removeFile

}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
