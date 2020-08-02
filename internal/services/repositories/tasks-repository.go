package repositories

import "headless-todo-file-service/internal/entities"

type TasksRepository interface {
	AddFileToTask(file entities.File, taskId string) error
}
