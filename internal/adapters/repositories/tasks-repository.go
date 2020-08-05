package repositories

import (
	"context"
	"encoding/json"
	"errors"
	kitnats "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services/repositories"
	"time"
)

const FileAddedSubjectName = "tasks.files.added"
const GetTaskByIdSubjectName = "tasks.getById"

type addFileToTaskRequest struct {
	TaskId   string `json:"taskId"`
	FileId   string `json:"fileId"`
	FileName string `json:"fileName"`
}

type getTaskByIdRequest struct {
	TaskId string `json:"taskId"`
}

type tasksRepositoryNats struct {
	sc stan.Conn
}

func (t *tasksRepositoryNats) GetTaskById(ctx context.Context, taskId string) (*entities.Task, error) {
	publisher := kitnats.NewPublisher(t.sc.NatsConn(), GetTaskByIdSubjectName, kitnats.EncodeJSONRequest, func(ctx context.Context, msg *nats.Msg) (response interface{}, err error) {
		var task entities.Task
		err = json.Unmarshal(msg.Data, &task)
		if err != nil {
			return nil, err
		}

		return task, nil
	})
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()
	errs := make(chan error)
	tasks := make(chan entities.Task)

	go func() {
		response, err := publisher.Endpoint()(ctx, getTaskByIdRequest{taskId})
		if err != nil {
			errs <- err
			return
		}
		task, ok := response.(entities.Task)
		if !ok {
			errs <- errors.New("wrong response structure")
			return
		}
		tasks <- task
	}()

	select {
	case task := <-tasks:
		return &task, nil
	case err := <-errs:
		return nil, err
	case <-ctx.Done():
		return nil, errors.New("could not verify the task's owner")
	}
}

func NewTasksRepositoryNats(sc stan.Conn) repositories.TasksRepository {
	return &tasksRepositoryNats{sc}
}

func (t *tasksRepositoryNats) AddFileToTask(file entities.File, taskId string) error {
	addFileToTaskRequest := &addFileToTaskRequest{taskId, file.ID.Hex(), file.Name}
	req, err := json.Marshal(addFileToTaskRequest)
	if err != nil {
		return err
	}
	err = t.sc.Publish(FileAddedSubjectName, req)
	if err != nil {
		return err
	}
	return nil
}
