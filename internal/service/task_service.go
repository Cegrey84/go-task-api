package service

import (
    "fmt"
    "github.com/Cegrey84/go-task-api/internal/models"
    "github.com/Cegrey84/go-task-api/internal/repository"
)

type TaskService struct {
    repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

func (s *TaskService) GetAllTasks() []models.Task {
    return s.repo.GetAll()
}

func (s *TaskService) GetTaskByID(id uint) (*models.Task, bool) {
    task, _ := s.repo.GetByID(id)
    if task == nil {
        return nil, false
    }
    return task, true
}

func (s *TaskService) CreateTask(text string) (*models.Task, error) {
    if text == "" {
        return nil, fmt.Errorf("text is required")
    }
    
    task := &models.Task{
        Text:   text,
        IsDone: false,
    }
    s.repo.Create(task)
    return task, nil
}

func (s *TaskService) UpdateTask(id uint, text string, isDone bool) (*models.Task, bool) {
    return s.repo.Update(id, text, isDone)
}

func (s *TaskService) DeleteTask(id uint) bool {
    return s.repo.Delete(id)
}