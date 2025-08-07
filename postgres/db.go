package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	var err error

	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

	db, err = sql.Open("postgres", dsn)

	if err != nil {
		return err
	}
	return db.Ping()
}

func GetOrderByID(orderUID string) (*Order, error) {
	order := Order{}

	err := db.QueryRow(`
        SELECT *
        FROM orders WHERE order_uid = $1`, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSig, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return nil, fmt.Errorf("orders not found: %w", err)
	}

	err = db.QueryRow(`
        SELECT *
        FROM delivery WHERE order_uid = $1`, orderUID).Scan(
		&order.Delivery.OrderUID, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email)
	if err != nil {
		return nil, fmt.Errorf("delivery not found: %w", err)
	}

	err = db.QueryRow(`
	    SELECT *
	    FROM payment WHERE order_uid = $1`, orderUID).Scan(
		&order.Payment.Transaction, &order.Payment.OrderUID, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank,
		&order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	rows, err := db.Query(`
	    SELECT *
	    FROM items WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, fmt.Errorf("items not found: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.OrderUID, &item.ChrtID, &item.TrackNumber,
			&item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func PrintDatabase(order Order) {
	orderJSON, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при сериализации:", err)
		return
	}

	fmt.Println("Order:")
	fmt.Println(string(orderJSON))
}

func CloseDB() {
	db.Close()
}
