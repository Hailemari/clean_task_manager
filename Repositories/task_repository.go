package repositories

import (
	"context"
	"errors"

	"github.com/Hailemari/clean_architecture_task_manager/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTaskRepository struct {
    collection *mongo.Collection
}

func NewMongoTaskRepository(collection *mongo.Collection) domain.TaskRepository{
    return &MongoTaskRepository{collection: collection}
}

func (r *MongoTaskRepository) GetTasks() ([]domain.Task, error) {
    var tasks []domain.Task
    cursor, err := r.collection.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var task domain.Task
        if err := cursor.Decode(&task); err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return tasks, nil
}

func (r *MongoTaskRepository) GetTaskByID(id primitive.ObjectID) (domain.Task, bool, error) {
    var task domain.Task
    err := r.collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&task)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return task, false, nil
        }
        return task, false, err
    }
    return task, true, nil
}

func (r *MongoTaskRepository) AddTask(task domain.Task) error {
    _, err := r.collection.InsertOne(context.TODO(), task)
    return err
}

func (r *MongoTaskRepository) UpdateTask(id primitive.ObjectID, task domain.Task) error {
    result := r.collection.FindOneAndUpdate(
        context.TODO(),
        bson.D{{Key: "_id", Value: id}},
        bson.D{{Key: "$set", Value: task}},
    )
    if result.Err() != nil {
        if result.Err() == mongo.ErrNoDocuments {
            return errors.New("task not found")
        }
        return result.Err()
    }
    return nil
}

func (r *MongoTaskRepository) DeleteTask(id primitive.ObjectID) error {
    result, err := r.collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}})
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        return mongo.ErrNoDocuments
    }
    return nil
}
