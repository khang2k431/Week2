package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func Init() {
	once.Do(func() {
		// build postgres DSN from env
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		name := os.Getenv("DB_NAME")

		if host != "" && port != "" && user != "" && pass != "" && name != "" {
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
				host, user, pass, name, port)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				DB = db
				log.Println("Connected to Postgres DB")
				return
			}
			log.Printf("Postgres connect failed: %v — fallback to sqlite", err)
		}

		// fallback to sqlite in-memory file (for dev)
		db, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to open sqlite:", err)
		}
		DB = db
		log.Println("Using fallback sqlite DB (dev.db)")
	})
}
