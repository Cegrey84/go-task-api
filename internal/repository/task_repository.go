package repository

import (
	"github.com/Cegrey84/go-task-api/internal/models"
	"sync"
	"time"
)

type TaskRepository struct {
	tasks  []models.Task
	mu     sync.RWMutex
	nextID uint
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks:  []models.Task{},
		nextID: 1,
	}
}

func (r *TaskRepository) GetAll() []models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.tasks
}

func (r *TaskRepository) GetByID(id uint) (*models.Task, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i, task := range r.tasks {
		if task.ID == id {
			return &task, i
		}
	}
	return nil, -1
}

func (r *TaskRepository) Create(task *models.Task) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	r.nextID++
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	r.tasks = append(r.tasks, *task)
}

func (r *TaskRepository) Update(id uint, text string, isDone bool) (*models.Task, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			if text != "" {
				r.tasks[i].Text = text
			}
			r.tasks[i].IsDone = isDone
			r.tasks[i].UpdatedAt = time.Now()
			return &r.tasks[i], true
		}
	}
	return nil, false
}

func (r *TaskRepository) Delete(id uint) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return true
		}
	}
	return false
}
