package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var task string

func getTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "hello, %s", task)
}

func postTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type requestBody struct {
		Task string `json:"task"`
	}

	var body requestBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	task = body.Task

	fmt.Fprintf(w, "Task saved: %s\n", body.Task)
}

func main() {
	http.HandleFunc("/", getTask)
	http.HandleFunc("/task", postTask)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
