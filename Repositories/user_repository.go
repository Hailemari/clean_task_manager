package repositories

import (
	"errors"
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/Hailemari/clean_architecture_task_manager/Domain"
)

type MongoUserRepository struct {
    collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
    return &MongoUserRepository{collection: collection}
}

func (r *MongoUserRepository) CreateUser(user *domain.User) error {
    existingUser := &domain.User{}
    err := r.collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(existingUser)
    if err == nil {
        return errors.New("user already exists")
    }

    count, err := r.collection.CountDocuments(context.TODO(), bson.M{})
    if err != nil {
        return err
    }
    if count == 0 {
        user.Role = "admin"
    } else {
        user.Role = "user"
    }

    _, err = r.collection.InsertOne(context.TODO(), user)
    return err
}

func (r *MongoUserRepository) GetUserByUsername(username string) (*domain.User, error) {
    user := &domain.User{}
    err := r.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(user)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *MongoUserRepository) PromoteUser(username string) error {
    var user bson.M
    err := r.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return errors.New("user not found")
        }
        return err
    }

    if userRole, ok := user["role"].(string); ok && userRole == "admin" {
        return errors.New("user is already an admin")
    }

    result, err := r.collection.UpdateOne(
        context.TODO(),
        bson.M{"username": username},
        bson.M{"$set": bson.M{"role": "admin"}},
    )

    if err != nil {
        return err
    }
    if result.ModifiedCount == 0 {
        return errors.New("user not found")
    }
    return nil
}