package endpoints

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"github.com/sony/gobreaker"
	"go.uber.org/dig"
	"headless-todo-file-service/internal/adapters/middlewares"
	"headless-todo-file-service/internal/entities"
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

type createFileResponse struct {
	File entities.File `json:"file,omitempty"`
	Err  error         `json:"error,omitempty"`
}

func (c createFileResponse) Error() error {
	return c.Err
}

func (c *createFileRequest) SetUserClaim(claim UserClaim) {
	c.UserClaim = claim
}

func makeCreateFileEndpoint(service services.FilesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*createFileRequest)
		file, err := service.Create(ctx, req.Name, req.UserClaim.ID, req.TaskId, req.File)
		if err != nil {
			return createFileResponse{}, err
		}
		return createFileResponse{
			File: *file,
		}, nil
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
	breaker := circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "Create-File",
		MaxRequests: 100,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}))
	taskEndpoint = breaker(taskEndpoint)
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
		GetDefaultHTTPOptions()...,
	)
}
