package request

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecordingTasksAndRetrievingThem(t *testing.T) {

	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := NewFileSystemTaskStore(database)
	assertNoError(t, err)

	server := NewTaskServer(store)

	//file_system_store_test.go
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemTaskStore(database)

		assertNoError(t, err)
	})

	t.Run("post a task and retrieve it", func(t *testing.T) {
		newTask := `{"id":1,"title":"A","description":"First Task"}`
		postRequest, _ := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(newTask))
		postResponse := httptest.NewRecorder()
		server.ServeHTTP(postResponse, postRequest)
		assertStatusCode(t, postResponse.Code, http.StatusAccepted)

		getRequest, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
		getResponse := httptest.NewRecorder()
		server.ServeHTTP(getResponse, getRequest)
		assertStatusCode(t, getResponse.Code, http.StatusOK)

		var got Task
		json.NewDecoder(getResponse.Body).Decode(&got)
		want := Task{ID: 1, Title: "A", Description: "First Task"}
		assertTask(t, got, want)
	})

	t.Run("post multiple tasks and retrieve all", func(t *testing.T) {
		newTask := `{"id":2,"title":"B","description":"Second Task"}`
		postRequest, _ := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(newTask))
		postResponse := httptest.NewRecorder()
		server.ServeHTTP(postResponse, postRequest)

		getAllRequest, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		getAllResponse := httptest.NewRecorder()
		server.ServeHTTP(getAllResponse, getAllRequest)
		assertStatusCode(t, getAllResponse.Code, http.StatusOK)

		var got []Task
		json.NewDecoder(getAllResponse.Body).Decode(&got)
		want := []Task{
			{ID: 1, Title: "A", Description: "First Task"},
			{ID: 2, Title: "B", Description: "Second Task"},
		}
		assertTasks(t, got, want)
	})

	t.Run("delete a task and verify its gone", func(t *testing.T) {
		deleteRequest, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
		deleteResponse := httptest.NewRecorder()
		server.ServeHTTP(deleteResponse, deleteRequest)
		assertStatusCode(t, deleteResponse.Code, http.StatusOK)

		getRequest, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
		getResponse := httptest.NewRecorder()
		server.ServeHTTP(getResponse, getRequest)
		assertStatusCode(t, getResponse.Code, http.StatusNotFound)
	})

	t.Run("update a task and verify the change", func(t *testing.T) {
		updatedTask := `{"title":"Updated B","description":"Updated Second Task"}`
		putRequest, _ := http.NewRequest(http.MethodPut, "/tasks/2", strings.NewReader(updatedTask))
		putResponse := httptest.NewRecorder()
		server.ServeHTTP(putResponse, putRequest)
		assertStatusCode(t, putResponse.Code, http.StatusOK)

		getRequest, _ := http.NewRequest(http.MethodGet, "/tasks/2", nil)
		getResponse := httptest.NewRecorder()
		server.ServeHTTP(getResponse, getRequest)

		var got Task
		json.NewDecoder(getResponse.Body).Decode(&got)
		want := Task{ID: 2, Title: "Updated B", Description: "Updated Second Task"}
		assertTask(t, got, want)
	})

}
