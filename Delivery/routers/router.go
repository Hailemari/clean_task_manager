package routers

import (
	"github.com/Hailemari/clean_architecture_task_manager/Domain"
	"github.com/Hailemari/clean_architecture_task_manager/Infrastructure"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the routes and middleware for the application
func SetupRouter(taskCtrl domain.TaskControllerInterface, userCtrl domain.UserControllerInterface) *gin.Engine {
    r := gin.Default()

    // Public routes
    r.POST("/register", userCtrl.CreateUser)
    r.POST("/login", userCtrl.LoginUser)

    // Protected routes
    auth := r.Group("/")
    auth.Use(infrastructure.AuthMiddleware())
    {
        auth.GET("/tasks", taskCtrl.GetTasks)
        auth.GET("/tasks/:id", taskCtrl.GetTask)

        // Admin-only routes
        admin := auth.Group("/")
        admin.Use(infrastructure.AdminMiddleware())
        {
            admin.POST("/tasks", taskCtrl.AddTask)
            admin.PUT("/tasks/:id", taskCtrl.UpdateTask)
            admin.DELETE("/tasks/:id", taskCtrl.DeleteTask)
            admin.POST("/promote/:username", userCtrl.PromoteUser)
        }
    }

    return r
}
