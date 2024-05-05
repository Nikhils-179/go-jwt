package routes

import (
	controller "github.com/Nikhils-179/go-jwt/controllers"
	middleware "github.com/Nikhils-179/go-jwt/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.POST("users/logout",controller.Logout())
}
