package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskFile struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name"`
}

type Task struct {
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	UserId      string              `json:"userId"`
	Files       []TaskFile          `json:"files"`
	Created     primitive.Timestamp `json:"created"`
}
