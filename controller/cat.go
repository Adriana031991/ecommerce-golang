package controller

import (
	"context"
	"ecommerce/golang/database"
	"ecommerce/golang/models"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}


func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		// productID, err := primitive.ObjectIDFromHex(productQueryID)
		// if err != nil {
		// 	log.Println(err)
		// 	c.AbortWithStatus(http.StatusInternalServerError)
		// 	return
		// }
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productQueryID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully Added to the cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("pid")
		if productQueryID == "" {
			log.Panicln("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		
		userQueryID := c.Query("uid")
		if userQueryID == "" {
			log.Panicln("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}


		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productQueryID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			
			c.IndentedJSON(200, "Successfully removed from cart")
			
					}
	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("uid")
		if userQueryID == "" {
			log.Println("userId is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid id"})
			c.Abort()
			return
		}

		usertID, err := primitive.ObjectIDFromHex(userQueryID)
		if err != nil {
			log.Println("aa",err)
		}
		
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcart models.User

		err = UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usertID}}).Decode(&filledcart)
		if err != nil {
			log.Println("bb",err)
			c.IndentedJSON(500, "not id found")
			return
		}

		filter_match := bson.D{{Key:"$match", Value: bson.D{primitive.E{Key:"_id", Value: usertID}}}}
		unwind := bson.D{{Key:"$unwind", Value: bson.D{primitive.E{Key:"path", Value:"$user_cart"}}}}
		grouping := bson.D{{Key:"$group", Value: bson.D{primitive.E{Key:"_id", Value:"$_id"}, {Key:"total", Value: bson.D{primitive.E{Key:"$sum", Value: "$user_cart.price"}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
		}

		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		// a := pointcursor.All(ctx, &listing)
		// fmt.Println("pt", a)
		// fmt.Println("pt2", listing)

		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledcart.UserCart)
		}
		ctx.Done()
	}
}

func(app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		userQueryID := c.Query("uid")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully placed the order")
	}
}

func(app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("pid")
		if productQueryID == "" {
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		
		userQueryID := c.Query("uid")
		if userQueryID == "" {
			log.Println("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println("-->",err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Successfully placed to order")
	}
}
