package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/grom"
)

var (
	DB *gorm.DB
	once sync.Once
)

func Init() {
	once.Do(func() {
		// build Postgres DSN from env
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_POST")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		name := os.Getenv("DB_NAME")

		if host != "" && port != "" && user != "" && pass != "" && name != "" {
			dsn := fmt.Sprint(
				"host=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
				host, user, pass, name, port,
			)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				// test connection
				if err := db.Exec("SELECT 1").Error; err == nil {
					DB = db
					log.Println(" Connected to Postgres DB")
					return
				}
				log.Printf(" Postgres connected but test query failed: %v", err)
			} else {
				log.Printf(" Postgres connect failed: %v, err)
			}
	   }

	   // fallback to sqlite
	   db, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
	   if err != nil {
		   log.Fatal(" Failed to open sqlite:", err)
		}
		DB = db
		log.Println(" Using fallback SQLite DB (dev.db)")
   })
}
