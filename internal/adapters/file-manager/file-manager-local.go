package file_manager

import (
	"bufio"
	"headless-todo-file-service/internal/adapters/repositories"
	"io"
	"os"
)

type FileManagerLocal struct {
	pathToLocalStorage string
}

func NewFileManagerLocal(pathToLocalStorage string) repositories.FileManager {
	return &FileManagerLocal{pathToLocalStorage: pathToLocalStorage}
}

func (f *FileManagerLocal) getFileFullName(id string) string {
	return f.pathToLocalStorage + "/" + id
}

func deleteFile(path string) error {
	return os.Remove(path)
}

func (f *FileManagerLocal) SaveFile(id string, file io.Reader) error {
	fileFullName := f.getFileFullName(id)
	createdFile, err := os.Create(fileFullName)
	if err != nil {
		_ = deleteFile(fileFullName)
		return err
	}

	reader := bufio.NewReader(file)
	_, err = reader.WriteTo(createdFile)
	if err != nil {
		_ = deleteFile(fileFullName)
		return err
	}
	return nil
}

func (f *FileManagerLocal) DeleteFile(id string) error {
	fileFullName := f.getFileFullName(id)
	return os.Remove(fileFullName)
}
