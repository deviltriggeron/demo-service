package main

import (
	"encoding/json"
	"log"
	"service/service"

	cache "cache/cache"
	kafka "kafka/kafka"
	model "model/model"
	db "postgres/postgres"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatal("Error init database: ", err)
	}

	defer db.CloseDB()
	kafka.ClearTopic("orders", 0)
	kafka.CreateTopic("orders", 1, 1)

	var orderCache = cache.NewCache()

	svc := service.NewOrderService(orderCache)

	order := model.Order{
		OrderUID: "a1234cd56ef78gh90xyz", TrackNumber: "WBILMTESTTRACK2", Entry: "WEB", Locale: "ru", InternalSig: "sig98765", CustomerID: "customer42",
		DeliveryService: "dhl", Shardkey: "3", SmID: 77, DateCreated: "2022-02-15T10:45:30Z", OofShard: "2",

		Delivery: model.Delivery{
			OrderUID: "a1234cd56ef78gh90xyz", Name: "Ivan Ivanov", Phone: "+79998887766", Zip: "101000", City: "Moscow", Address: "Tverskaya Street, 12",
			Region: "Moscow region", Email: "ivanov@example.com",
		},

		Payment: model.Payment{
			OrderUID: "a1234cd56ef78gh90xyz", Transaction: "a1234cd56ef78gh90xyz", RequestID: "req-555", Currency: "RUB", Provider: "sberpay", Amount: 5599,
			PaymentDT: 1644912330, Bank: "T-Bank", DeliveryCost: 299, GoodsTotal: 5300, CustomFee: 0,
		},

		Items: []model.Item{
			{OrderUID: "a1234cd56ef78gh90xyz", ChrtID: 88776655, TrackNumber: "WBILMTESTTRACK2", Price: 2650, Rid: "rid9988776655", Name: "Wireless Headphones", Sale: 10,
				Size: "M", TotalPrice: 2385, NmID: 7654321, Brand: "Sony", Status: 201,
			},
			{OrderUID: "a1234cd56ef78gh90xyz", ChrtID: 22334455, TrackNumber: "WBILMTESTTRACK2", Price: 2700, Rid: "rid22334455test", Name: "Smart Watch", Sale: 15,
				Size: "L", TotalPrice: 2295, NmID: 8765432, Brand: "Iphone", Status: 202,
			},
		},
	}

	orders, err := svc.GetAll()
	if err != nil {
		log.Println("Get error:", err)
	} else {
		res, err := json.MarshalIndent(orders, "", " ")
		if err != nil {
			log.Println("Indent error:", err)
			return
		}
		println("From cache", string(res))
	}
	kafka.ProduceOrder("orders", "UPDATE", order)
	// kafka.ProduceOrder("orders", "INSERT", order)
	// kafka.ProduceOrder("orders", "SELECT", model.Order{})

	go kafka.ConsumeOrder("orders", svc)
	select {}
}

/* init db

produce
consume
init cache
add to cache

print from cache or db
*/

// model.Order{
// 	OrderUID: "b563feb7b2b84b6test", TrackNumber: "WBILMTESTTRACK", Entry: "WBIL", Locale: "en", InternalSig: "", CustomerID: "test",
// 	DeliveryService: "meest", Shardkey: "9", SmID: 99, DateCreated: "2021-11-26T06:22:19Z", OofShard: "1",

// 	Delivery: model.Delivery{
// 		Name: "Test Testov", Phone: "+9720000000", Zip: "2639809", City: "Kiryat Mozkin", Address: "Ploshad Mira 15",
// 		Region: "Kraiot", Email: "test@gmail.com",
// 	},

// 	Payment: model.Payment{
// 		Transaction: "b563feb7b2b84b6test", RequestID: "", Currency: "USD", Provider: "wbpay", Amount: 1817,
// 		PaymentDT: 1637907727, Bank: "alpha", DeliveryCost: 1500, GoodsTotal: 317, CustomFee: 0,
// 	},

// 	Items: []model.Item{
// 		{ChrtID: 9934930, TrackNumber: "WBILMTESTTRACK", Price: 453, Rid: "ab4219087a764ae0btest", Name: "Mascaras", Sale: 30,
// 			Size: "0", TotalPrice: 317, NmID: 2389212, Brand: "Vivienne Sabo", Status: 202},
// 	},
// }
