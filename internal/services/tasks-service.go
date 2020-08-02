package services

import (
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services/repositories"
)

type TasksService interface {
	AddFileToTask(file entities.File, taskId string) error
}

type tasksService struct {
	tasksRepository repositories.TasksRepository
}

func NewTasksService(tasksRepository repositories.TasksRepository) TasksService {
	return &tasksService{tasksRepository: tasksRepository}
}

func (t *tasksService) AddFileToTask(file entities.File, taskId string) error {
	return t.tasksRepository.AddFileToTask(file, taskId)
}
