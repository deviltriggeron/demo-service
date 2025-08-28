package main

import (
	cache "cache/cache"
	kafka "kafka/kafka"
	"log"
	db "postgres/postgres"
	server "server/server"
	"service/service"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatal("Error init database: ", err)
	}

	defer db.CloseDB()
	// kafka.ClearTopic("orders", 0)
	kafka.CreateTopic("orders", 1, 1)

	var orderCache = cache.NewCache()

	svc := service.NewOrderService(orderCache)

	go kafka.ConsumeOrder("orders", svc)

	httpSrv := server.NewServer(svc)
	go httpSrv.Run()

	select {}
}

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

// model.Order{
// 	OrderUID:       "a12345b6c7890test",
// 	TrackNumber:    "WBNEWTRACK2025",
// 	Entry:          "WBX",
// 	Locale:         "ru",
// 	InternalSig:    "sig123",
// 	CustomerID:     "customer42",
// 	DeliveryService:"dhl",
// 	Shardkey:       "5",
// 	SmID:           77,
// 	DateCreated:    "2025-08-27T12:45:00Z",
// 	OofShard:       "3",

// 	Delivery: model.Delivery{
// 		Name:    "Ivan Ivanov",
// 		Phone:   "+79001234567",
// 		Zip:     "101000",
// 		City:    "Moscow",
// 		Address: "Tverskaya 10",
// 		Region:  "Moscow Region",
// 		Email:   "ivanov@example.com",
// 	},

// 	Payment: model.Payment{
// 		Transaction:  "a12345b6c7890test",
// 		RequestID:    "req-555",
// 		Currency:     "EUR",
// 		Provider:     "paypal",
// 		Amount:       2599,
// 		PaymentDT:    1693123456,
// 		Bank:         "sberbank",
// 		DeliveryCost: 499,
// 		GoodsTotal:   2100,
// 		CustomFee:    0,
// 	},

// 	Items: []model.Item{
// 		{
// 			ChrtID:      11223344,
// 			TrackNumber: "WBNEWTRACK2025",
// 			Price:       700,
// 			Rid:         "rid-xyz-001",
// 			Name:        "Wireless Headphones",
// 			Sale:        15,
// 			Size:        "M",
// 			TotalPrice:  595,
// 			NmID:        445566,
// 			Brand:       "Sony",
// 			Status:      201,
// 		},
// 	},
// }

// model.Order{
// 	OrderUID:        "b98765c4d3210test",
// 	TrackNumber:     "WBTESTTRACK2025",
// 	Entry:           "WBY",
// 	Locale:          "en",
// 	InternalSig:     "sig987",
// 	CustomerID:      "customer77",
// 	DeliveryService: "fedex",
// 	Shardkey:        "8",
// 	SmID:            99,
// 	DateCreated:     "2025-08-28T09:30:00Z",
// 	OofShard:        "1",

// 	Delivery: model.Delivery{
// 		Name:    "John Smith",
// 		Phone:   "+12025550123",
// 		Zip:     "10001",
// 		City:    "New York",
// 		Address: "5th Avenue 101",
// 		Region:  "NY",
// 		Email:   "john.smith@example.com",
// 	},

// 	Payment: model.Payment{
// 		Transaction:  "b98765c4d3210test",
// 		RequestID:    "req-777",
// 		Currency:     "USD",
// 		Provider:     "stripe",
// 		Amount:       4999,
// 		PaymentDT:    1694123456,
// 		Bank:         "chase",
// 		DeliveryCost: 999,
// 		GoodsTotal:   4000,
// 		CustomFee:    0,
// 	},

// 	Items: []model.Item{
// 		{
// 			ChrtID:      55667788,
// 			TrackNumber: "WBTESTTRACK2025",
// 			Price:       2000,
// 			Rid:         "rid-abc-002",
// 			Name:        "Gaming Mouse",
// 			Sale:        10,
// 			Size:        "L",
// 			TotalPrice:  1800,
// 			NmID:        778899,
// 			Brand:       "Logitech",
// 			Status:      202,
// 		},
// 		{
// 			ChrtID:      99887766,
// 			TrackNumber: "WBTESTTRACK2025",
// 			Price:       2500,
// 			Rid:         "rid-def-003",
// 			Name:        "Mechanical Keyboard",
// 			Sale:        20,
// 			Size:        "XL",
// 			TotalPrice:  2000,
// 			NmID:        889900,
// 			Brand:       "Razer",
// 			Status:      202,
// 		},
// 	},
// }
