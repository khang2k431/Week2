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
	// Load .env (if exists)
	_ = godotenv.Load()

	// Init DB (Postgres or fallback sqlite)
	config.Init()

	if config.DB == nil {
		log.Fatal(" Database connection failed")
	}

	// Auto-migrate models
	if err := config.DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatal(" Migration failed:", err)
	}

	// Create default admin if not exists
	createAdminIfNotExist()

	// setup gin 
	r := gin.Default()

	// Global middlewares
	r.Use(middlewares.RateLimitMiddleware())

	// Public routes
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	// Protected routes
	auth := r.Group("/api")
	auth.Use(middlewares.JWTAuthMiddleware())
	{
		auth.POST("/tasks", controllers.CreateTask)
		auth.GET("/tasks", controllers.ListTasks)
		auth.GET("/tasks/:id", controllers.GetTask)
		auth.PUT("/tasks/:id", controllers.UpdateTask)
		auth.DELETE("/tasks/:id", controllers.DeleteTask)
	}

	// Run server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf(" Server running on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal(" Server failed:", err)
	}

	// Create default admin account if not exists
	func createAdminIfNotExist() {
		var count int64
		config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
		if count == 0 {
			pwd := "admin123"
			hashed, _ := bcrybt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
			admin := models.User{
				Username: "admin",
				Email:    "admin@example.com",
				Password: string(hashed),
				Role:     "admin",
		}
		config.DB.Create(&admin)
		log.Prinln(" Created default admin account (email=admin@example.com, password=admin123)")
		}
	}
