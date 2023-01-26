package worker

import (
	"database/sql"
	"food_delivery_mysql/models"
	"food_delivery_mysql/utils"
	"log"
	"sync"
	"time"
)

func StartWorker() {
	wg := sync.WaitGroup{}
	go DeliveryBoy()
	wg.Add(1)

	go Customer()
	wg.Add(1)

	go Restaurant()
	wg.Add(1)

	wg.Wait()
}

func DeliveryBoy() {

}

func Customer() {

}

func Restaurant() {

}

func handleOrders() {
	db, err := utils.Database()
	if err != nil {
		log.Println("error occurred while getting database", err)
		return
	}

	ticker := time.NewTicker(time.Second * 1)
	for range ticker.C {
		query := `select * from foodorder where cookassigned = false`
		accounts, err := utils.GetAccounts(db, query)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, account := range accounts {
			err = notifyRestaurant(db, account.OrderID, account.RestaurantName)
			if err != nil {
				log.Println("error occurred while updating restaurant")
				continue
			}
		}
	}
}

func cookFood(db *sql.DB, customerid int) (err error) {
	query := `update foodorder set cookassigned = true where cookedandready is false, customerid = ?`

	result, err := db.Exec(query, customerid)
	if err != nil {
		log.Println("error occurred while getting food orders")
		return
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil || rowsaffected == 0 {
		log.Println("failed to update food order", rowsaffected, err)
		return
	}
	utils.MyPrint("started cooking")
	time.Sleep(time.Second * 30)

	query = `update foodorder set cookedandready = true where customerid = ?, completestatus = false`

	result, err = db.Exec(query, customerid)
	if err != nil {
		log.Println("error occurred while getting food orders")
		return
	}
	rowsaffected, err = result.RowsAffected()
	if err != nil || rowsaffected == 0 {
		log.Println("failed to update ready food order", rowsaffected, err)
		return
	}
	return
}

func notifyRestaurant(db *sql.DB, orderid int, restaurantname string) (err error) {
	query := `update restaurant set orderid = ? where name = ?`
	result, err := db.Exec(query, orderid, restaurantname)
	if err != nil {
		return err
	}
	if rowsaffected, err := result.RowsAffected(); rowsaffected == 0 || err != nil {
		log.Println("no resords were updated", err, rowsaffected)
		return err
	}
	return nil
}

func notifyDeliveryBoy(db *sql.DB, customerid int) {

}

func getOrder() {

}

func getFoodOrders(db *sql.DB) (err error) {
	query := `select * from foodorder where completestatus = false && boyid is not null`

	rows, err := db.Query(query)
	if err != nil {
		log.Println("error occurred while getting food orders")
		return
	}
	orders := []models.FoodOrder{}
	for rows.Next() {
		singleOrder := models.FoodOrder{}
		err = rows.Scan(&singleOrder.ID, &singleOrder.FoodName, &singleOrder.Quantity, &singleOrder.RestaurantName,
			&singleOrder.CustomerID, &singleOrder.PlacedTime, &singleOrder.CookAssigned,
			&singleOrder.BoyID, &singleOrder.CookedAndReady, &singleOrder.CompleteStatus)

		if err != nil {
			log.Println("scanning error occurred", err)
			return
		}

		orders = append(orders, singleOrder)
	}
	log.Println("started cooking ")
	return
}
