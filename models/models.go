package models

import (
	"time"
)

type FoodOrder struct {
	ID             int    `bson:"id" json:"id"`
	FoodName       string `bson:"foodname" json:"food_name"`
	Quantity       int    `bson:"quantity" json:"quantity"`
	RestaurantName string `bson:"restaurantname" json:"restaurant_name"`
	CustomerID     int    `bson:"customerid" json:"customer_id"`
	PlacedTime     string `bson:"placedtime" json:"placed_time"`
	CookAssigned   bool   `bson:"cookassigned" json:"cookassigned"`
	BoyID          int    `bson:"boyid,omitempty" json:"boy_id"`
	CookedAndReady bool   `bson:"cookedandready" json:"cooked_and_ready"`
	CompleteStatus bool   `bson:"completestatus" json:"complete_status"`
}

// CREATE TABLE foodorder (id int NOT NULL AUTO_INCREMENT, foodname VARCHAR(20), quantity int, restaurantname VARCHAR(20), customerid int, address VARCHAR(20), placedtime DATE, cookassigned BOOLEAN, boyid int, cookedandready BOOLEAN, completestatus BOOLEAN, PRIMARY KEY (id));

type Restaurant struct {
	ID          int    `bson:"id,omitempty" json:"id"`
	Name        string `bson:"name" json:"name"`
	Address     string `bson:"address" json:"address"`
	FoodOrderID int    `bson:"foodorderid,omitempty" json:"food_order_id"`
}

// CREATE TABLE restaurant (id int NOT NULL AUTO_INCREMENT, name VARCHAR(20),  address VARCHAR(20), foodorderid int, PRIMARY KEY (id));

type DeliveryBoy struct {
	ID              int    `bson:"id,omitempty" json:"id"`
	Name            string `bson:"name" json:"name"`
	BusyStatus      bool   `bson:"busystatus" json:"busystatus"`
	OrderID         int    `bson:"orderid" json:"order_id"`
	CurrentLocation string `bson:"currentlocation" json:"current_location"`
}

// CREATE TABLE deliveryboy (id int NOT NULL AUTO_INCREMENT, name VARCHAR(20),  busystatus BOOLEAN, currentorderid int, currentlocation VARCHAR(20), PRIMARY KEY (id));

type CustomerAccount struct {
	ID             int       `bson:"id" json:"id"`
	Name           string    `bson:"name" json:"name"`
	Phone          string    `bson:"phone" json:"phone"`
	Address        string    `bson:"address" json:"address"`
	OrderID        int       `bson:"orderid" json:"order_id"`
	OrderPlaced    bool      `bson:"orderplaced" json:"order_placed"`
	RestaurantName string    `bson:"restaurantname" json:"restaurant_name"`
	PlacedTime     time.Time `bson:"placedtime" json:"placed_time"`
	ReceiveTime    time.Time `bson:"receivetime" json:"receive_time"`
	BoyName        string    `bson:"boyname" json:"boy_name"`
	InProcess      bool      `bson:"inprocess" json:"inprocess"`
}

// CREATE TABLE customer (id int NOT NULL AUTO_INCREMENT, name VARCHAR(20), phone VARCHAR(20), address VARCHAR(20), currentorderid int, orderplaced BOOLEAN, placedtime DATE, receivetime DATE, boyname VARCHAR(20), inprocess BOOLEAN, PRIMARY KEY (id));
