package main

import (
	"food_delivery_mysql/controller"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	controller.Controller()
}
