package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

/*
Updating the functions to have error returns because in postgres errors are normal

Giving (int, error) to addtask specifically because now postgres generates the id,
the client doesn't supplu it anymore.
*/
type TaskStore interface {
	GetAllTasks() (List, error)
	AddTask(task Task) (int, error)
	DeleteTask(id int) (bool, error)
	UpdateTask(id int, task Task) (bool, error)
}

type TaskServer struct {
	store TaskStore
	http.Handler
}

func NewTaskServer(store TaskStore) *TaskServer {

	t := &TaskServer{
		store: store,
	}

	router := http.NewServeMux()
	router.Handle("/tasks", http.HandlerFunc(t.handleTaskCollection))
	router.Handle("/tasks/", http.HandlerFunc(t.handleTaskItem))
	t.Handler = router

	return t
}

func (t *TaskServer) handleTaskCollection(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		t.returnAllTasks(w)
	case http.MethodPost:
		t.handlePostRequest(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (t *TaskServer) handleTaskItem(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		t.handleGetRequest(w, r)
	case http.MethodDelete:
		t.handleDeleteRequest(w, r)
	case http.MethodPut:
		t.handlePutRequest(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (t *TaskServer) returnAllTasks(w http.ResponseWriter) {

	tasks, err := t.store.GetAllTasks()

	if err != nil {
		http.Error(w, "failed to retreive tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (t *TaskServer) returnTask(w http.ResponseWriter, id int) {

	tasks, err := t.store.GetAllTasks()

	if err != nil {
		http.Error(w, "failed to retreive tasks", http.StatusInternalServerError)
		return
	}

	task, found := tasks.Find(id)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (t *TaskServer) handlePostRequest(w http.ResponseWriter, r *http.Request) {

	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	newId, err := t.store.AddTask(task)
	if err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	task.ID = newId
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated) // changed from stausAccepted to statusCreated.

	json.NewEncoder(w).Encode(task)

}

func (t *TaskServer) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if path[2] == "" {
		t.returnAllTasks(w)
		return
	}

	id, err := strconv.Atoi(path[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	t.returnTask(w, id)
}

func (t *TaskServer) handleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	idstr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ok, err := t.store.DeleteTask(id)
	if err != nil {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (t *TaskServer) handlePutRequest(w http.ResponseWriter, r *http.Request) {
	idstr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var task Task

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := t.store.UpdateTask(id, task)
	if err != nil {
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
