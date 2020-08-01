package repositories

import (
	"context"
	"headless-todo-file-service/internal/entities"
	"io"
)

type FilesRepository interface {
	Create(context.Context, string, string, io.Reader) (*entities.File, error)
}
