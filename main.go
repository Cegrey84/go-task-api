package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Text   string `json:"text"`
	IsDone bool   `json:"is_done"`
}

var db *gorm.DB

func main() {

	var err error
	db, err = gorm.Open(sqlite.Open("file:tasks.db?_pragma=busy_timeout(5000)"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = db.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		var tasks []Task
		db.Find(&tasks)
		json.NewEncoder(w).Encode(tasks)
	}).Methods("GET")

	r.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		var task Task
		json.NewDecoder(r.Body).Decode(&task)

		result := db.Create(&task)
		if result.Error != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
	}).Methods("POST")

	r.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		var task Task
		db.First(&task, id)

		var updates map[string]interface{}
		json.NewDecoder(r.Body).Decode(&updates)

		if text, ok := updates["text"].(string); ok {
			task.Text = text
		}
		if isDone, ok := updates["is_done"].(bool); ok {
			task.IsDone = isDone
		}

		db.Save(&task)
		json.NewEncoder(w).Encode(task)
	}).Methods("PATCH")

	r.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		var task Task
		db.First(&task, id)
		task.IsDone = true
		db.Save(&task)

		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Database: tasks.db")
	fmt.Println("Using modernc.org/sqlite driver")
	http.ListenAndServe(":8080", r)
}
