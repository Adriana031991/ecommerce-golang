package main

import (
	"log"
	"os"

	"ecommerce/golang/controller"
	"ecommerce/golang/database"
	"ecommerce/golang/middleware"
	"ecommerce/golang/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := controller.NewApplication(
		database.ProductData(database.Client, "Products"),
		database.UserData(database.Client, "Users"))
	
	router := gin.New()
	router.Use(gin.Logger())
	
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/listcart", controller.GetItemFromCart())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	router.GET("/addtocart", app.AddToCart())
	router.PUT("/removeitem", app.RemoveItem())
	router.POST("/addaddress", controller.AddAddress())
	router.PUT("/edithomeaddress", controller.EditHomeAddress())
	router.PUT("/editworkaddress", controller.EditWorkAddress())
	router.PUT("/deleteaddresses", controller.DeleteAddress())

	log.Fatal(router.Run(":" + port))
}