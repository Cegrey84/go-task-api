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

type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var tasks = []Task{}
var nextID uint = 1

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSONResponse{
		Success: true,
		Data:    tasks,
	})
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JSONResponse{
			Success: false,
			Error:   "Неверный JSON формат",
		})
		return
	}

	if task.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JSONResponse{
			Success: false,
			Error:   "Поле 'text' обязательно для заполнения",
		})
		return
	}

	task.ID = nextID
	nextID++
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	tasks = append(tasks, task)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSONResponse{
		Success: true,
		Message: "Задача создана успешно",
		Data:    task,
	})
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JSONResponse{
			Success: false,
			Error:   "Неверный ID задачи",
		})
		return
	}

	var foundTask *Task
	for i := range tasks {
		if tasks[i].ID == uint(id) {
			foundTask = &tasks[i]
			break
		}
	}

	if foundTask == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSONResponse{
			Success: false,
			Error:   "Задача не найдена",
		})
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JSONResponse{
			Success: false,
			Error:   "Неверный JSON формат",
		})
		return
	}

	if text, ok := updates["text"].(string); ok && text != "" {
		foundTask.Text = text
	}

	if isDone, ok := updates["is_done"].(bool); ok {
		foundTask.IsDone = isDone
	}

	foundTask.UpdatedAt = time.Now()

	json.NewEncoder(w).Encode(JSONResponse{
		Success: true,
		Message: "Задача обновлена успешно",
		Data:    foundTask,
	})
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JSONResponse{
			Success: false,
			Error:   "Неверный ID задачи",
		})
		return
	}

	for i := range tasks {
		if tasks[i].ID == uint(id) {
			tasks[i].IsDone = true
			tasks[i].UpdatedAt = time.Now()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(JSONResponse{
		Success: false,
		Error:   "Задача не найдена",
	})
}

func setupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PATCH")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Task API Server (In-Memory Version)",
			"version": "1.0.0",
			"docs":    "Available endpoints: GET/POST /tasks, PATCH/DELETE /tasks/{id}",
		})
	}).Methods("GET")

	return r
}

func printStartupInfo() {
	fmt.Println(" Task API Server Started! (In-Memory Version)")
	fmt.Println(" URL: http://localhost:8080")
	fmt.Println("")
	fmt.Println(" Available endpoints:")
	fmt.Println(" GET    /tasks          - Get all tasks")
	fmt.Println(" POST   /tasks          - Create new task")
	fmt.Println(" PATCH  /tasks/{id}     - Update task")
	fmt.Println(" DELETE /tasks/{id}     - Mark task as done")
	fmt.Println("")
	fmt.Println(" Example POST JSON:")
	fmt.Println(`   {"text": "Buy milk", "is_done": false}`)
	fmt.Println("")
	fmt.Println(" Ready for testing!")
	fmt.Println("")
}

func main() {

	r := setupRoutes()

	printStartupInfo()

	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(" Ошибка запуска сервера:", err)
	}
}
