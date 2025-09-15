package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;size:100" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:200" json:"email"`
	Password  string         `json:"-"`
	Role      string         `gorm:"size:20" json:"role"` // admin or user
	Tasks     []Task         `json:"tasks"`
}

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"size:255" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Category    string         `gorm:"size:100" json:"category"`
	DueDate     *time.Time     `json:"due_date"`
	Completed   bool           `json:"completed"`
	OwnerID     uint           `json:"owner_id"`
	Owner       User           `json:"owner" gorm:"foreignKey:OwnerID"`
}
