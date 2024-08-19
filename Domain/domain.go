package domain

import (
    "time"
    "errors"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title       string             `json:"title"`
    Description string             `json:"description"`
    DueDate     time.Time          `json:"due_date"`
    Status      string             `json:"status"`
}

type User struct {
    ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Username string             `json:"username"`
    Password string             `json:"password"`
    Role     string             `json:"role"`
}

var AllowedStatuses = []string{"pending", "in-progress", "completed"}

func (t *Task) Validate() error {
    if t.ID == primitive.NilObjectID {
        return errors.New("task ID cannot be empty")
    }
    if t.Title == "" {
        return errors.New("task title cannot be empty")
    }
    if t.DueDate.IsZero() {
        return errors.New("task due date cannot be empty")
    }
    if t.Status == "" {
        return errors.New("task status cannot be empty")
    }
    if !t.isValidStatus() {
        return errors.New("invalid task status. Allowed statuses are: pending, in-progress, completed")
    }
    return nil
}

func (t *Task) isValidStatus() bool {
    for _, allowedStatus := range AllowedStatuses {
        if t.Status == allowedStatus {
            return true
        }
    }
    return false
}

func (u *User) ValidateUser() error {
    if u.Username == "" {
        return errors.New("username cannot be empty")
    }
    if u.Password == "" {
        return errors.New("password cannot be empty")
    }
    return nil
}

func (u *User) HashPassword() error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

func (u *User) ComparePassword(password string) error {
    return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type TaskRepository interface {
    GetTasks() ([]Task, error)
    GetTaskByID(id primitive.ObjectID) (Task, bool, error)
    AddTask(task Task) error
    UpdateTask(id primitive.ObjectID, task Task) error
    DeleteTask(id primitive.ObjectID) error
}

type UserRepository interface {
    CreateUser(user *User) error
    GetUserByUsername(username string) (*User, error)
    PromoteUser(username string) error
}

type TaskUseCaseInterface interface {
    GetTasks() ([]Task, error)
    GetTask(id primitive.ObjectID) (Task, bool, error)
    AddTask(task Task) error
    UpdateTask(id primitive.ObjectID, task Task) error
    DeleteTask(id primitive.ObjectID) error
}

type UserUseCaseInterface interface {
    CreateUser(user *User) error
    GetUserByUsername(username string) (*User, error)
    PromoteUser(username string) error
}

type TaskControllerInterface interface {
    GetTasks(ctx *gin.Context)
    GetTask(ctx *gin.Context)
    AddTask(ctx *gin.Context)
    UpdateTask(ctx *gin.Context)
    DeleteTask(ctx *gin.Context)
}

type UserControllerInterface interface {
    CreateUser(ctx *gin.Context)
    LoginUser(ctx *gin.Context)
    PromoteUser(ctx *gin.Context)
}