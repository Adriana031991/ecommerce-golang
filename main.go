package main

import (
	"log"
	"os"

	"ecommerce/golang/controllers"
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

	app := controllers.NewApplication(
		database.productData(database.Client, "products"),
		database.UserData(database.Client, "users"))
	
	router := gin.New()
	router.Use(gin.Logger())
	
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.Get("/addtocart", app.AddToCart())
	router.Get("/removeitem", app.RemoveItem())
	router.Get("/cartcheckout", app.BuyFromCart())
	router.Get("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}