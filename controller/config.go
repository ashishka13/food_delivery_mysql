package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"food_delivery_mysql/models"
	"food_delivery_mysql/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func createAccount(w http.ResponseWriter, r *http.Request) {
	account := models.CustomerAccount{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		log.Println("decode error occurred", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("%+v", account)
	db, err := utils.Database()
	if err != nil {
		log.Println("error occurred while connecting database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO `customer` (`name`, `phone`, `address`) VALUES (?, ?, ?)"
	insertResult, err := db.ExecContext(context.Background(), query, account.Name, account.Phone, account.Address)
	if err != nil {
		log.Printf("impossible insert teacher: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	lastinsertID, err := insertResult.LastInsertId()
	if err != nil {
		log.Printf("impossible to retrieve last inserted id: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("inserted id: %d", lastinsertID)
	w.WriteHeader(http.StatusOK)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	log.Printf("%+v", name)

	db, err := utils.Database()
	if err != nil {
		log.Println("error occurred while connecting database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := db.Exec("DELETE FROM `customer` WHERE `name` = ?", name)
	if err != nil {
		log.Println("error occurred while deleting record", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("occurred while processing the result", err)
		return
	}
	if rows == 0 {
		log.Println("record not found to delete")
		bytemsg, err := json.Marshal("record not found with that name")
		if err != nil {
			log.Println("error occurred while decoding", err)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(bytemsg)
		return
	}

	log.Println("record deleted successfully")
	w.WriteHeader(http.StatusOK)
}

func placeOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerid"]
	intcustomerID, err := strconv.Atoi(customerID)
	if err != nil {
		log.Println("error occurred in string conversion", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	order := models.FoodOrder{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&order)
	if err != nil {
		log.Println("decode error occurred", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	order.PlacedTime = time.Now().Format("2006-01-02 15:04:05")
	log.Printf("%+v", order)

	db, err := utils.Database()
	if err != nil {
		log.Println("error occurred while connecting database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	query := `INSERT INTO foodorder (foodname, quantity, restaurantname, customerid, placedtime, cookassigned,
	boyid, cookedandready, completestatus)	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	insertResult, err := db.ExecContext(context.Background(), query, order.FoodName, order.Quantity, order.RestaurantName,
		intcustomerID, order.PlacedTime, order.CookAssigned, order.BoyID, order.CookedAndReady,
		order.CompleteStatus)
	if err != nil {
		log.Printf("error occurred while inserting order: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	orderid, err := insertResult.LastInsertId()
	if err != nil {
		log.Printf("impossible to retrieve last inserted id: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("update customer set currentorderid = ?, orderplaced = ?, placedtime = ? where id = ?",
		orderid, true, order.PlacedTime, intcustomerID)
	if err != nil {
		log.Println("error occurred while updating record", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytemsg, err := json.Marshal("record successfully updated")
	if err != nil {
		log.Println("error occurred while decoding", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = notifyRestaurant(w, db, int(orderid), order.RestaurantName)
	if err != nil {
		log.Println("error occurred while updating restaurant", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytemsg)
}

func notifyRestaurant(w http.ResponseWriter, db *sql.DB, orderid int, restaurantname string) (err error) {
	query := `INSERT INTO restaurant (name, address, foodorderid) VALUES (?, ?, ?)`

	_, err = db.ExecContext(context.Background(), query, restaurantname, restaurantname, orderid)
	if err != nil {
		log.Printf("error occurred while inserting order: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return nil
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	db, err := utils.Database()
	if err != nil {
		log.Println("error occurred while connecting database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rows, err := db.Query("select * from foodorder")
	if err != nil {
		return
	}
	allOrders := []models.FoodOrder{}
	for rows.Next() {
		singleOrder := models.FoodOrder{}
		err := rows.Scan(&singleOrder.ID, &singleOrder.FoodName, &singleOrder.Quantity,
			&singleOrder.RestaurantName, &singleOrder.CustomerID, &singleOrder.PlacedTime,
			&singleOrder.CookAssigned, &singleOrder.BoyID,
			&singleOrder.CookedAndReady, &singleOrder.CompleteStatus)

		if err != nil {
			log.Println("error occurred while decoding rows", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		allOrders = append(allOrders, singleOrder)
	}

	for _, val := range allOrders {
		log.Println(val)
	}

	w.WriteHeader(http.StatusOK)
}

func InsertCustomer(db *sql.DB, w http.ResponseWriter, account models.CustomerAccount) (lastinsertID int64, err error) {
	query := "INSERT INTO `customer` (`name`, `phone`, `address`) VALUES (?, ?, ?)"
	insertResult, err := db.ExecContext(context.Background(), query, account.Name, account.Phone, account.Address)
	if err != nil {
		log.Printf("impossible insert teacher: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	lastinsertID, err = insertResult.LastInsertId()
	if err != nil {
		log.Printf("impossible to retrieve last inserted id: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func DeleteCustomer(db *sql.DB, w http.ResponseWriter, id int, name string) (err error) {
	query := `DELETE FROM customer WHERE name = ?`
	result, err := db.Exec(query, name)
	if err != nil {
		log.Println("error occurred while deleting record", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("occurred while processing the result", err)
		return
	}
	if rows == 0 {
		var bytemsg []byte
		log.Println("record not found to delete")
		bytemsg, err = json.Marshal("record not found with that name")
		if err != nil {
			log.Println("error occurred while decoding", err)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(bytemsg)
		return
	}
	return
}
