package repositories

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services/repositories"
)

const FileAddedSubjectName = "tasks.files.added"

type addFileToTaskRequest struct {
	TaskId   string `json:"taskId"`
	FileId   string `json:"fileId"`
	FileName string `json:"fileName"`
}

type tasksRepositoryNats struct {
	nc *nats.Conn
}

func NewTasksRepositoryNats(nc *nats.Conn) repositories.TasksRepository {
	return &tasksRepositoryNats{nc: nc}
}

func (t *tasksRepositoryNats) AddFileToTask(file entities.File, taskId string) error {
	addFileToTaskRequest := &addFileToTaskRequest{taskId, file.ID.Hex(), file.Name}
	req, err := json.Marshal(addFileToTaskRequest)
	if err != nil {
		return err
	}
	err = t.nc.Publish(FileAddedSubjectName, req)
	if err != nil {
		return err
	}
	return nil
}
