package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/Hailemari/clean_architecture_task_manager/Domain"
    "github.com/Hailemari/clean_architecture_task_manager/Usecases"
    "github.com/Hailemari/clean_architecture_task_manager/Infrastructure"
)

type Controller struct {
    taskUseCase *usecases.TaskUseCase
    userUseCase *usecases.UserUseCase
}

func NewController(taskUC *usecases.TaskUseCase, userUC *usecases.UserUseCase) *Controller {
    return &Controller{
        taskUseCase: taskUC,
        userUseCase: userUC,
    }
}

func (c *Controller) GetTasks(ctx *gin.Context) {
    tasks, err := c.taskUseCase.GetTasks()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (c *Controller) GetTask(ctx *gin.Context) {
    id := ctx.Param("id")
    task, found, err := c.taskUseCase.GetTask(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
        return
    }
    if !found {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
        return
    }
    ctx.JSON(http.StatusOK, task)
}

func (c *Controller) AddTask(ctx *gin.Context) {
    var newTask domain.Task
    if err := ctx.ShouldBindJSON(&newTask); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.taskUseCase.AddTask(newTask); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func (c *Controller) UpdateTask(ctx *gin.Context) {
    id := ctx.Param("id")
    var updatedTask domain.Task
    if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.taskUseCase.UpdateTask(id, updatedTask); err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (c *Controller) DeleteTask(ctx *gin.Context) {
    id := ctx.Param("id")
    err := c.taskUseCase.DeleteTask(id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func (c *Controller) CreateUser(ctx *gin.Context) {
    var newUser domain.User
    if err := ctx.ShouldBindJSON(&newUser); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.userUseCase.CreateUser(&newUser); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (c *Controller) LoginUser(ctx *gin.Context) {
    var loginData struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := ctx.ShouldBindJSON(&loginData); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := c.userUseCase.GetUserByUsername(loginData.Username)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := user.ComparePassword(loginData.Password); err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := infrastructure.GenerateToken(user)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *Controller) PromoteUser(ctx *gin.Context) {
    username := ctx.Param("username")
    if err := c.userUseCase.PromoteUser(username); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}