package middlewares

import (
	"context"
	"github.com/go-kit/kit/log"
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services"
	"time"
)

type LoggerMiddleware struct {
	Logger log.Logger
	Next   services.TasksService
}

func (l *LoggerMiddleware) Create(ctx context.Context, name, description, userId string) (output *entities.Task, err error) {
	defer func(begin time.Time) {
		_ = l.Logger.Log(
			"method", "Create",
			"name", name,
			"description", description,
			"userId", userId,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return l.Next.Create(ctx, name, description, userId)
}
