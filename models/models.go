package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	First_name      *string            `json:"name" validate:"required,min=2,max=30" bson:"name"`
	Last_Name       *string            `json:"last_name"  validate:"required,min=2,max=30" bson:"last_name"`
	Password        *string            `json:"password" validate:"required,min=6" bson:"password" `
	Email           *string            `json:"email"  validate:"required" bson:"email"`
	Phone           *string            `json:"phone"  validate:"required" bson:"phone"`
	Token           *string            `json:"token" bson:"token"`
	Refresh_token   *string            `json:"refresh_token" bson:"refresh_token"`
	Created_At      time.Time          `json:"created_at" bson:"created_at"`
	Updated_At      time.Time          `json:"updated_at" bson:"updated_at"`
	User_id         string             `json:"user_id" bson:"user_id"`
	UserCart        []ProductUser      `json:"user_cart" bson:"user_cart"`
	Address_Details []Address          `json:"address_details" bson:"address_details"`
	Order_Status    []Order            `json:"order_status" bson:"order_status"`
}

type Product struct {
	Product_ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        int                `json:"price" bson:"price"`
	Rating       *uint8             `json:"rating" bson:"rating"`
	Image        *string            `json:"image" bson:"image"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        int                `json:"price" bson:"price"`
	Rating       *uint8             `json:"rating" bson:"rating"`
	Image        *string            `json:"image" bson:"image"`
	
}

type Address struct {
	Address_Id primitive.ObjectID `json:"id" bson:"_id"`
	House      *string            `json:"house" bson:"house"`
	Street     *string            `json:"street" bson:"street"`
	City       *string            `json:"city" bson:"city"`
	PinCode    *string            `json:"pincode" bson:"pincode"`
}

type Order struct {
	Order_Id       primitive.ObjectID `json:"id" bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_cart" bson:"order_cart"`
	Ordered_At     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price          int                `json:"price" bson:"price"`
	Discount       *int               `json:"discount" bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool
	Cod     bool
}
