package controllers

import (
	"net/http"

	"github.com/Hailemari/clean_architecture_task_manager/Domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
    useCase domain.TaskUseCaseInterface
}

func NewTaskController(useCase domain.TaskUseCaseInterface) domain.TaskControllerInterface {
    return &TaskController{useCase: useCase}
}

func (c *TaskController) GetTasks(ctx *gin.Context) {
    tasks, err := c.useCase.GetTasks()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, tasks)
}

func (c *TaskController) GetTask(ctx *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
    task, found, err := c.useCase.GetTask(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if !found {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        return
    }
    ctx.JSON(http.StatusOK, task)
}

func (c *TaskController) AddTask(ctx *gin.Context) {
    var task domain.Task
    if err := ctx.ShouldBindJSON(&task); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := c.useCase.AddTask(task); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"message": "task added"})
}

func (c *TaskController) UpdateTask(ctx *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
    var task domain.Task
    if err := ctx.ShouldBindJSON(&task); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := c.useCase.UpdateTask(id, task); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "task updated"})
}

func (c *TaskController) DeleteTask(ctx *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
    if err := c.useCase.DeleteTask(id); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

type UserController struct {
    useCase domain.UserUseCaseInterface
}

func NewUserController(useCase domain.UserUseCaseInterface) domain.UserControllerInterface {
    return &UserController{useCase: useCase}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
    var user domain.User
    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := c.useCase.CreateUser(&user); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (c *UserController) LoginUser(ctx *gin.Context) {
    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := c.useCase.GetUserByUsername(input.Username)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    if err := user.ComparePassword(input.Password); err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // Generate token here (omitted for brevity)

    ctx.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (c *UserController) PromoteUser(ctx *gin.Context) {
    username := ctx.Param("username")
    if err := c.useCase.PromoteUser(username); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "user promoted"})
}
