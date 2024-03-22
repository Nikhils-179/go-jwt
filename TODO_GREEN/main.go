package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	// Initialize GORM
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	// Auto-migrate the schema
	if err := db.AutoMigrate(&Todo{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// Create a Gin router
	router := gin.Default()

	// Handlers for various endpoints
	router.POST("/todos", func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		if err := db.Create(&todo).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create Todo"})
			return
		}

		c.JSON(200, todo)
	})

	router.GET("/todos", func(c *gin.Context) {
		var todos []Todo
		if err := db.Find(&todos).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch Todos"})
			return
		}

		c.JSON(200, todos)
	})

	// Implement other CRUD operations handlers as needed

	// Run the server on port 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}



