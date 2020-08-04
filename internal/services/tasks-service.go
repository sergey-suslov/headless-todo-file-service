package services

import (
	"context"
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services/repositories"
)

type TasksService interface {
	AddFileToTask(file entities.File, taskId string) error
	GetTaskById(ctx context.Context, taskId string) (*entities.Task, error)
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

func (t *tasksService) GetTaskById(ctx context.Context, taskId string) (*entities.Task, error) {
	task, err := t.tasksRepository.GetTaskById(ctx, taskId)
	if err != nil {
		return nil, err
	}
	return task, nil
}
