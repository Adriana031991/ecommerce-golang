package database

import (
	"context"
	"ecommerce/golang/models"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't decode the product")
	ErrUserIdIsNotValid = errors.New("this user is not a valid")
	ErrProductIdIsNotValid = errors.New("this product is not a valid")
	ErrCantUpdateUser = errors.New("can't add this product to the cart")
	ErrCantRemoveItem = errors.New("can't remove this item from the cart")
	ErrCantGetItem = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem = errors.New("can't update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productID string, userID string) error {
	var productcart models.ProductUser

	_pid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		log.Println(err)
		return ErrProductIdIsNotValid
	}

	err = prodCollection.FindOne(ctx, bson.M{"_id": _pid}).Decode(&productcart)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	// err = searchfromdb.All(ctx, &productcart)
	// if err != nil {
	// 	log.Println(err)
	// 	return ErrCantDecodeProducts
	// }

	_uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: _uid}}
	NewUpdate := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "user_cart", Value: productcart}}}}
	// update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "user_cart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}
	_, err = userCollection.UpdateOne(ctx, filter, NewUpdate)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID string, userID string) error {
	_pid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		log.Println(err)
		return ErrProductIdIsNotValid
	}
	
	
	_uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	fmt.Println("-->", _pid)
	


	filter := bson.D{primitive.E{Key: "_id", Value: _uid}}
	update := bson.D{{Key: "$pull", Value: bson.D{primitive.E{Key: "user_cart", Value: bson.D{primitive.E{Key: "_id", Value: _pid}}}}}}
	// update := bson.M{"$pull": bson.M{"user_cart": bson.M{"_id": pid}}}
	fmt.Println("aa",update.Map())
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItem
	}
	return nil

}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	//fetch the cart of the user

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_Id = primitive.NewObjectID()
	orderCart.Ordered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.Cod = true
	
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$user_cart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$user_cart.price"}}}}}}
	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	
	ctx.Done()
	if err != nil {
		panic(err)
	}
	
	//find the cart total
	var getUserCart []bson.M
	if err = currentResults.All(ctx, &getUserCart); err != nil {
		fmt.Println(err)
		panic(err)
	}
	
	var total_price int32
	for _, user_item := range getUserCart {
		price := user_item["total"]
		total_price = price.(int32)
	}
	orderCart.Price = int(total_price)
	
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "order_status", Value: orderCart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}
	
	//create an order with the items
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"order_status.$[].order_cart": bson.M{"$each": getCartItems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	//empty up the cart
	userCart_empty := make([]models.ProductUser, 0)
	filtered := bson.D{primitive.E{Key: "_id", Value: id}}
	updated := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "user_cart", Value: userCart_empty}}}}
	_, err = userCollection.UpdateOne(ctx, filtered, updated)
	if err != nil {
		return ErrCantBuyCartItem

	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, UserID string) error {
	id, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var product_details models.ProductUser
	var orders_detail models.Order
	
	orders_detail.Order_Id = primitive.NewObjectID()
	orders_detail.Ordered_At = time.Now()
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method.Cod = true

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
	}
	orders_detail.Price = product_details.Price
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "order_status", Value: orders_detail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"order_status.$[].order_cart": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}