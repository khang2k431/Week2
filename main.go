package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"Week2/config"
	"Week2/controllers"
	"Week2/middlewares"
	"Week2/models"
)

func main() {
	// load .env (nếu có)
	_ = godotenv.Load()

	// init DB (Postgres if available, fallback to sqlite)
	config.Init()

	// auto-migrate models
	if err := config.DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatal("migration failed:", err)
	}

	// create default admin if not exists
	createAdminIfNotExist()

	r := gin.Default()

	// global middlewares
	r.Use(middlewares.RateLimitMiddleware())

	// public
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	// protected routes
	auth := r.Group("/api")
	auth.Use(middlewares.JWTAuthMiddleware())
	{
		auth.POST("/tasks", controllers.CreateTask)
		auth.GET("/tasks", controllers.ListTasks)
		auth.GET("/tasks/:id", controllers.GetTask)
		auth.PUT("/tasks/:id", controllers.UpdateTask)
		auth.DELETE("/tasks/:id", controllers.DeleteTask)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}
}

func createAdminIfNotExist() {
	var count int64
	config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		pwd := "admin123"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		admin := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(hashed),
			Role:     "admin",
		}
		config.DB.Create(&admin)
		log.Printf("Created default admin: email=admin@example.com password=%s", pwd)
	}
}
