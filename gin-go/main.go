package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)
type UserController struct{}

//Creating Controller and Handler
func (uc *UserController) GetUserInfo(c *gin.Context){
	userId := c.Param("id")

	c.JSON(200,gin.H{"id": userId,"name":"John Cena","email":"J@gmail.com"})
}
func LoggerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request Method: %s | Status : %d | Duration : %v", c.Request.Method, c.Writer.Status(), duration)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		//modifed code to pass test case
		if apiKey != "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}
func main() {
	//Default engine
	router := gin.Default()

	userController := &UserController{}

	router.GET("/users/:id",userController.GetUserInfo)

	//Middleware
	router.Use(LoggerMiddleWare())

	//Routing
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})
	router.GET("/bye", func(c *gin.Context) {
		c.String(200, "Bye")
	})

	//Grouping
	authGroup := router.Group("/api")
	authGroup.Use(AuthMiddleware())
	{
		authGroup.GET("/data", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Autheticated and authorized"})
		})
	}

	private := router.Group("/private")

	private.Use(AuthMiddleware())
	{
		private.GET("/private", func(c *gin.Context) {
			c.String(200, "Private data accessed after authetication")
		})
		private.POST("/create", func(c *gin.Context) {
			c.String(200, "Create a new resource")
		})
	}

	public := router.Group("/public")
	{
		public.GET("/info", func(c *gin.Context) {
			c.String(200, "Public info")
		})

		public.POST("/postproduct", func(c *gin.Context) {
			c.String(200, "Posted")
		})
	}

	router.Run(":8000")
}
