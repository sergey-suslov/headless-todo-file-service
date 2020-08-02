package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services/repositories"
	"io"
)

const FilesCollection = "files"

type FilesRepositoryMongo struct {
	db *mongo.Database
	fm FileManager
}

func NewFilesRepositoryMongo(db *mongo.Database, fm FileManager) repositories.FilesRepository {
	return &FilesRepositoryMongo{db, fm}
}

func (r *FilesRepositoryMongo) deleteFile(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.Collection(FilesCollection).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *FilesRepositoryMongo) Create(ctx context.Context, name, userId string, file io.Reader) (*entities.File, error) {
	task := entities.NewFile(name, userId)
	result, err := r.db.Collection(FilesCollection).InsertOne(ctx, bson.M{"name": task.Name, "userId": task.UserId, "created": task.Created})
	if err != nil {
		return nil, err
	}
	task.ID = result.InsertedID.(primitive.ObjectID)

	err = r.fm.SaveFile(task.ID.Hex(), file)
	if err != nil {
		_ = r.deleteFile(ctx, task.ID)
		return nil, err
	}

	return &task, nil
}
