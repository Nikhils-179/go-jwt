package routes

import(
	controller "github.com/Nikhils-179/go-jwt/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("users/signup",controller.Signup())
	incomingRoutes.GET("users/login",controller.Login())
}