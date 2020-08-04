package repositories

import (
	"context"
	"headless-todo-file-service/internal/entities"
)

type TasksRepository interface {
	AddFileToTask(file entities.File, taskId string) error
	GetTaskById(ctx context.Context, taskId string) (*entities.Task, error)
}
