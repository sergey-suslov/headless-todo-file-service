package repositories

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
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
	sc stan.Conn
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
