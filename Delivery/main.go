package main

import (
	"os"
	"log"
	"context"

	"github.com/joho/godotenv"
	"github.com/Hailemari/clean_architecture_task_manager/Usecases"
	"github.com/Hailemari/clean_architecture_task_manager/Repositories"
	"github.com/Hailemari/clean_architecture_task_manager/Delivery/routers"
	"github.com/Hailemari/clean_architecture_task_manager/Delivery/controllers"
	
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Printf("Warning: Error loading .env file: %v", err)
    }

    // Initialize the database connection
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        log.Fatal("MONGODB_URI environment variable is not set")
    }

    client, err := connectDB(mongoURI)
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
    defer client.Disconnect(context.TODO())

    // Initialize repositories
    taskRepo := repositories.NewMongoTaskRepository(client.Database("taskDB").Collection("tasks"))
    userRepo := repositories.NewMongoUserRepository(client.Database("taskDB").Collection("users"))

    // Initialize use cases
    taskUC := usecases.NewTaskUseCase(taskRepo)
    userUC := usecases.NewUserUseCase(userRepo)

    // Initialize controller
    ctrl := controllers.NewController(taskUC, userUC)

    // Set up router
    r := routers.SetupRouter(ctrl)

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    log.Printf("Server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Could not start server: %v", err)
    }
}

func connectDB(mongoURI string) (*mongo.Client, error) {
    clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure the connection is established
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}