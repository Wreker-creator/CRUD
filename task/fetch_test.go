package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type StubTaskStore struct {
	tasks []Task
}

func (s *StubTaskStore) UpdateTask(id int, task Task) bool {

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i] = task
			return true
		}
	}
	return false

}

func (s *StubTaskStore) DeleteTask(id int) bool {

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return true
		}
	}
	return false
}

func (s *StubTaskStore) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

func (s *StubTaskStore) GetAllTasks() List {
	return s.tasks
}

func TestSingularTaskFunctions(t *testing.T) {

	store := &StubTaskStore{
		tasks: []Task{
			{ID: 1, Title: "A", Description: "First Task"},
			{ID: 2, Title: "B", Description: "Second Task"},
		},
	}
	server := NewTaskServer(store)

	t.Run("return task with specific task id", func(t *testing.T) {

		id := 1
		request := newGetTaskRequest(id)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got Task
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Errorf("unable to decode response body: %v", err)
		}

		want := Task{ID: 1, Title: "A", Description: "First Task"}
		assertTask(t, got, want)
		assertStatusCode(t, response.Code, http.StatusOK)

	})

	t.Run("return empty for no task", func(t *testing.T) {
		id := 3
		request := newGetTaskRequest(id)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("it returns accepted on Post", func(t *testing.T) {

		newTask := `{"id":3,"title":"C","description":"Third Task"}`

		request, _ := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(newTask))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusAccepted)
		if len(store.tasks) != 3 {
			t.Errorf("not added to the current database")
		}

		want := Task{ID: 3, Title: "C", Description: "Third Task"}
		assertTask(t, store.tasks[2], want)

	})

}

func TestReturnAllTasks(t *testing.T) {

	store := &StubTaskStore{
		tasks: []Task{
			{ID: 1, Title: "A", Description: "First Task"},
			{ID: 2, Title: "B", Description: "Second Task"},
			{ID: 3, Title: "C", Description: "Third Task"},
		},
	}

	server := NewTaskServer(store)

	t.Run("check if all tasks are returned", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		respone := httptest.NewRecorder()
		server.ServeHTTP(respone, request)

		var got []Task
		err := json.NewDecoder(respone.Body).Decode(&got)
		if err != nil {
			t.Errorf("unable to decode response body: %v", err)
		}

		want := []Task{
			{ID: 1, Title: "A", Description: "First Task"},
			{ID: 2, Title: "B", Description: "Second Task"},
			{ID: 3, Title: "C", Description: "Third Task"},
		}

		assertStatusCode(t, respone.Code, http.StatusOK)
		assertTasks(t, got, want)

	})

}

func TestDeleteTask(t *testing.T) {

	t.Run("check if the task was deleted", func(t *testing.T) {

		store := &StubTaskStore{
			tasks: []Task{
				{ID: 1, Title: "A", Description: "First Task"},
				{ID: 2, Title: "B", Description: "Second Task"},
				{ID: 3, Title: "C", Description: "Third Task"},
			},
		}

		server := NewTaskServer(store)
		id := 1

		request := newDeleteRequest(id)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusOK)

		want := []Task{
			{ID: 2, Title: "B", Description: "Second Task"},
			{ID: 3, Title: "C", Description: "Third Task"},
		}

		assertTasks(t, store.tasks, want)

	})

	t.Run("check if an error was raised for incorrect id", func(t *testing.T) {

		store := &StubTaskStore{
			tasks: []Task{
				{ID: 1, Title: "A", Description: "First Task"},
				{ID: 2, Title: "B", Description: "Second Task"},
				{ID: 3, Title: "C", Description: "Third Task"},
			},
		}

		server := NewTaskServer(store)
		id := 99

		request := newDeleteRequest(id)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusNotFound)
		want := []Task{
			{ID: 1, Title: "A", Description: "First Task"},
			{ID: 2, Title: "B", Description: "Second Task"},
			{ID: 3, Title: "C", Description: "Third Task"},
		}
		assertTasks(t, store.tasks, want)

	})

}

func TestUpdateTask(t *testing.T) {

	t.Run("check if the error is raised for non-existent id", func(t *testing.T) {
		store := &StubTaskStore{
			tasks: []Task{
				{ID: 1, Title: "A", Description: "First Task"},
				{ID: 2, Title: "B", Description: "Second Task"},
				{ID: 3, Title: "C", Description: "Third Task"},
			},
		}

		server := NewTaskServer(store)
		id := 4
		newTask := `{"id":3,"title":"Updated C","description":"Updated Third Task"}`
		request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", id), strings.NewReader(newTask))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusNotFound)

	})

	t.Run("check if the task is updated or not", func(t *testing.T) {

		store := &StubTaskStore{
			tasks: []Task{
				{ID: 1, Title: "A", Description: "First Task"},
				{ID: 2, Title: "B", Description: "Second Task"},
				{ID: 3, Title: "C", Description: "Third Task"},
			},
		}

		server := NewTaskServer(store)
		id := 3
		newTask := `{"id":3,"title":"Updated C","description":"Updated Third Task"}`
		request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", id), strings.NewReader(newTask))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []Task{
			{ID: 1, Title: "A", Description: "First Task"},
			{ID: 2, Title: "B", Description: "Second Task"},
			{ID: 3, Title: "Updated C", Description: "Updated Third Task"},
		}

		assertStatusCode(t, response.Code, http.StatusOK)
		assertTasks(t, store.GetAllTasks(), want)

	})

}

func newDeleteRequest(id int) *http.Request {
	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", id), nil)
	return request
}

func assertTask(t testing.TB, got, want Task) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func newGetTaskRequest(id int) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", id), nil)
	return req
}

func assertTasks(t testing.TB, got, want []Task) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d tasks, want %d", len(got), len(want))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf(
				"task mismatch at index %d\nGot:  %+v\nWant: %+v",
				i,
				got[i],
				want[i],
			)
		}
	}
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
