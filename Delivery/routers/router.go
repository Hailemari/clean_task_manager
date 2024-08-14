package routers

import (
    "github.com/gin-gonic/gin"
    "github.com/Hailemari/clean_architecture_task_manager/Infrastructure"
    "github.com/Hailemari/clean_architecture_task_manager/Delivery/controllers"
)

func SetupRouter(ctrl *controllers.Controller) *gin.Engine {
    r := gin.Default()

    // Public routes
    r.POST("/register", ctrl.CreateUser)
    r.POST("/login", ctrl.LoginUser)

    // Protected routes
    auth := r.Group("/")
    auth.Use(infrastructure.AuthMiddleware())
    {
        auth.GET("/tasks", ctrl.GetTasks)
        auth.GET("/tasks/:id", ctrl.GetTask)

        // Admin-only routes
        admin := auth.Group("/")
        admin.Use(infrastructure.AdminMiddleware())
        {
            admin.POST("/tasks", ctrl.AddTask)
            admin.PUT("/tasks/:id", ctrl.UpdateTask)
            admin.DELETE("/tasks/:id", ctrl.DeleteTask)
            admin.POST("/promote/:username", ctrl.PromoteUser)
        }
    }

    return r
}