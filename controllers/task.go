package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"Week2/config"
	"Week2/models"
	"Week2/utils"
)

type createTaskInput struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	DueDate     *string `json:"due_date"` // RFC3339 optional
}

type updateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	DueDate     *string `json:"due_date"`
	Completed   *bool   `json:"completed"`
}

func CreateTask(c *gin.Context) {
	claims := c.MustGet("claims").(*utils.Claims)

	var in createTaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var due *time.Time
	if in.DueDate != nil && *in.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *in.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date, use RFC3339"})
			return
		}
		due = &t
	}

	task := models.Task{
		Title:       in.Title,
		Description: in.Description,
		Category:    in.Category,
		DueDate:     due,
		OwnerID:     claims.UserID,
	}
	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func ListTasks(c *gin.Context) {
	var tasks []models.Task
	// simple pagination
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}
	config.DB.Preload("Owner").Limit(pageSize).Offset((page-1)*pageSize).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := config.DB.Preload("Owner").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	claims := c.MustGet("claims").(*utils.Claims)
	id := c.Param("id")

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	// only admin or owner
	if claims.Role != "admin" && task.OwnerID != claims.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	var in updateTaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if in.Title != nil {
		task.Title = *in.Title
	}
	if in.Description != nil {
		task.Description = *in.Description
	}
	if in.Category != nil {
		task.Category = *in.Category
	}
	if in.Completed != nil {
		task.Completed = *in.Completed
	}
	if in.DueDate != nil {
		if *in.DueDate == "" {
			task.DueDate = nil
		} else {
			t, err := time.Parse(time.RFC3339, *in.DueDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date"})
				return
			}
			task.DueDate = &t
		}
	}

	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	claims := c.MustGet("claims").(*utils.Claims)
	id := c.Param("id")

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if claims.Role != "admin" && task.OwnerID != claims.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}
	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
