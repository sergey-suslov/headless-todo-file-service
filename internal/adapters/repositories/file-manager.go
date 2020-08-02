package repositories

import "io"

type FileManager interface {
	SaveFile(id string, file io.Reader) error
	DeleteFile(id string) error
}
