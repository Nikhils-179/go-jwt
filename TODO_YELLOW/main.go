package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Todo struct {
	ID          int64    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	// Creates a gin engine responsible for HTTP requests, routing, middleware, and generating HTTP response
	router := gin.Default()

	// Open SQLite database connection
	db, err := sql.Open("sqlite", "todo.db")
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	// Ensure Todo table exists
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS todo (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			description TEXT
		);
	`); err != nil {
		log.Fatal("failed to create todo table:", err)
	}

	// Handlers for various endpoints
	router.POST("/todos", func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		result, err := db.Exec("INSERT INTO todo (title, description) VALUES (?, ?)", todo.Title, todo.Description)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create Todo"})
			return
		}

		todo.ID, _ = result.LastInsertId()
		c.JSON(200, todo)
	})

	router.GET("/todos", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, description FROM todo")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch Todos"})
			return
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next() {
			var todo Todo
			if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
				c.JSON(500, gin.H{"error": "Failed to scan Todo"})
				return
			}
			todos = append(todos, todo)
		}

		c.JSON(200, todos)
	})

	// Implement other CRUD operations handlers as needed

	// Run the server on port 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
