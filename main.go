package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Task struct {
	ID        uint      `json:"id"`
	Text      string    `json:"text"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var tasks = []Task{}
var nextID uint = 1

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	if task.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Text is required"})
		return
	}

	task.ID = nextID
	nextID++
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	tasks = append(tasks, task)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	var taskIndex = -1
	for i, task := range tasks {
		if task.ID == uint(id) {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	if text, ok := updates["text"].(string); ok && text != "" {
		tasks[taskIndex].Text = text
	}

	if isDone, ok := updates["is_done"].(bool); ok {
		tasks[taskIndex].IsDone = isDone
	}

	tasks[taskIndex].UpdatedAt = time.Now()
	json.NewEncoder(w).Encode(tasks[taskIndex])
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var taskIndex = -1
	for i, task := range tasks {
		if task.ID == uint(id) {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PATCH")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Task API",
			"docs":    "GET/POST /tasks, PATCH/DELETE /tasks/{id}",
		})
	})

	fmt.Println(" Server: http://localhost:8080")
	fmt.Println(" Endpoints: GET,POST /tasks  PATCH,DELETE /tasks/{id}")
	log.Fatal(http.ListenAndServe(":8080", r))
}
