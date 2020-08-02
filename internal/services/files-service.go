package services

import (
	"context"
	"errors"
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services/repositories"
	"io"
)

type FilesService interface {
	Create(ctx context.Context, name, userId, tasksId string, file io.Reader) (*entities.File, error)
}

type filesService struct {
	filesRepository repositories.FilesRepository
	tasksRepository repositories.TasksRepository
}

func NewFilesServiceService(filesRepository repositories.FilesRepository, tasksRepository repositories.TasksRepository) FilesService {
	return &filesService{filesRepository, tasksRepository}
}

func (service *filesService) Create(ctx context.Context, name, userId, tasksId string, file io.Reader) (*entities.File, error) {
	if tasksId == "" {
		return nil, errors.New("taskId must be present")
	}
	if name == "" {
		return nil, errors.New("name must be present")
	}
	if userId == "" {
		return nil, errors.New("userId must be present")
	}
	createdFile, err := service.filesRepository.Create(ctx, name, userId, file)
	if err != nil {
		return nil, err
	}
	err = service.tasksRepository.AddFileToTask(*createdFile, tasksId)
	if err != nil {
		// TODO delete file
		return nil, err
	}

	return createdFile, nil
}
