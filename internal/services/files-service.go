package services

import (
	"context"
	"errors"
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services/repositories"
	"io"
)

type FilesService interface {
	Create(ctx context.Context, name, userId string, file io.Reader) (*entities.File, error)
}

type filesService struct {
	tasksRepository repositories.FilesRepository
}

func NewTasksService(tasksRepository repositories.FilesRepository) FilesService {
	return &filesService{tasksRepository}
}

func (service *filesService) Create(ctx context.Context, name, userId string, file io.Reader) (*entities.File, error) {
	if name == "" {
		return nil, errors.New("name must be present")
	}
	if userId == "" {
		return nil, errors.New("userId must be present")
	}
	return service.tasksRepository.Create(ctx, name, userId, file)
}
