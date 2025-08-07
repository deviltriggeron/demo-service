package main

import (
	"log"
	db "postgres/postgres"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatal("DB error: ", err)
	}

	order, err := db.GetOrderByID("b563feb7b2b84b6test")

	if err != nil {
		log.Fatal("Get error: ", err)
	}

	db.PrintDatabase(*order)

	db.CloseDB()
}
