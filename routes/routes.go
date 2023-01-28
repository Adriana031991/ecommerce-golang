package routes

import (
	"ecommerce/golang/controllers"

	"github.com/gin-gonic/gin"
)


func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signUp", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addProduct", controllers.AddProduct())
	incomingRoutes.GET("/users/productView", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
}