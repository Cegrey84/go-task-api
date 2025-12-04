package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var tasks []Task
var nextID = 1

func getTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func postTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	newTask.ID = nextID
	nextID++
	tasks = append(tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func patchTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Only PATCH allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			var updates map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			if text, ok := updates["text"].(string); ok {
				tasks[i].Text = text
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getTasks).Methods("GET")
	r.HandleFunc("/task", postTask).Methods("POST")
	r.HandleFunc("/task/{id}", patchTask).Methods("PATCH")
	r.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}