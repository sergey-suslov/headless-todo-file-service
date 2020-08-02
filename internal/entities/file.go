package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type File struct {
	ID      primitive.ObjectID  `json:"_id" bson:"_id"`
	Name    string              `json:"name"`
	UserId  string              `json:"userId"`
	Created primitive.Timestamp `json:"created"`
}

func NewFile(name, userId string) File {
	return File{Name: name, UserId: userId, Created: primitive.Timestamp{T: uint32(time.Now().Unix())}}
}
