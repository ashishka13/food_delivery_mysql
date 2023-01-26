package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Controller ...
func Controller() {
	r := mux.NewRouter()

	r.HandleFunc("/", welcome).Methods("GET", "PUT", "POST", "DELETE")
	r.HandleFunc("/customer/account", createAccount).Methods("POST")
	r.HandleFunc("/customer/account/{name}", deleteAccount).Methods("DELETE")

	r.HandleFunc("/orders/{customerid}", placeOrder).Methods("POST")
	r.HandleFunc("/orders", getOrders).Methods("GET")

	log.Println("listening..")
	http.ListenAndServe(":8080", r)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>welcome to food delivery service, choose restaurant and place your order<h1>"))
}
