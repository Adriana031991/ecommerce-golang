package routes

import (
	"ecommerce/golang/controller"

	"github.com/gin-gonic/gin"
)


func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signUp", controller.SignUp())
	incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.POST("/admin/addProduct", controller.ProductViewerAdmin())
	incomingRoutes.GET("/users/productView", controller.SearchProduct())
	incomingRoutes.GET("/users/search", controller.SearchProductByQuery())
}