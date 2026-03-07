package request

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

type TaskStore interface {
	GetAllTasks() List
	AddTask(task Task)
	DeleteTask(id int) bool
	UpdateTask(id int, task Task) bool
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

	// for rturning all tasks and handling post request

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
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(t.store.GetAllTasks())
}

func (t *TaskServer) returnTask(w http.ResponseWriter, id int) {
	task, check := t.store.GetAllTasks().Find(id)
	w.Header().Set("content-type", "application/json")
	if check == false {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (t *TaskServer) handlePostRequest(w http.ResponseWriter, r *http.Request) {

	var task Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t.store.AddTask(task)
	w.WriteHeader(http.StatusAccepted)

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
	ok := t.store.DeleteTask(id)
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

	task.ID = id
	ok := t.store.UpdateTask(id, task)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
