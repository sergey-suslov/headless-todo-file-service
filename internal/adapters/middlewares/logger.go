package middlewares

import (
	"context"
	"github.com/go-kit/kit/log"
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services"
	"io"
	"time"
)

type LoggerMiddleware struct {
	Logger log.Logger
	Next   services.FilesService
}

func (l *LoggerMiddleware) Create(ctx context.Context, name, userId, taskId string, file io.Reader) (output *entities.File, err error) {
	defer func(begin time.Time) {
		_ = l.Logger.Log(
			"method", "Create",
			"name", name,
			"userId", userId,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return l.Next.Create(ctx, name, userId, taskId, file)
}
