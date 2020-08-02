package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"headless-todo-file-service/internal/adapters/middlewares"
	"headless-todo-file-service/internal/services"
	"log"
	"mime/multipart"
	"net/http"
)

type createFileRequest struct {
	UserClaim
	Name   string
	TaskId string
	File   multipart.File
}

func (c *createFileRequest) SetUserClaim(claim UserClaim) {
	c.UserClaim = claim
}

func makeCreateFileEndpoint(service services.FilesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*createFileRequest)
		task, err := service.Create(ctx, req.Name, req.UserClaim.ID, req.File)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}

func CreateFileHandler(c *dig.Container) http.Handler {
	var service services.FilesService
	if err := c.Invoke(func(s services.FilesService) {
		service = s
	}); err != nil {
		log.Fatal(err)
	}

	var logger kitlog.Logger

	if err := c.Invoke(func(log kitlog.Logger) {
		logger = log
	}); err != nil {
		log.Fatal(err)
	}

	service = &middlewares.LoggerMiddleware{Logger: kitlog.With(logger), Next: service}
	taskEndpoint := makeCreateFileEndpoint(service)

	return httptransport.NewServer(
		taskEndpoint,
		DefaultRequestDecoder(func(r *http.Request) (UserClaimable, error) {
			defer r.Body.Close()
			fileName := r.FormValue("name")
			taskId := r.FormValue("taskId")
			taskFile, _, err := r.FormFile("file")
			if err != nil && err != http.ErrMissingFile {
				return nil, errors.Wrap(err, "manifest file")
			}
			return &createFileRequest{
				Name:   fileName,
				File:   taskFile,
				TaskId: taskId,
			}, nil
		}),
		DefaultRequestEncoder,
	)
}
