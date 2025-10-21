package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"demo-service/internal/config"
	"demo-service/internal/handler"
	"demo-service/internal/infrastructure/cache"
	"demo-service/internal/infrastructure/kafka"
	"demo-service/internal/infrastructure/postgres"
	"demo-service/internal/router"
	"demo-service/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg := config.LoadConfigDB()

	storage, err := postgres.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error init database: %v", err)
	}

	orderCache := cache.NewCache(10*time.Minute, 15*time.Minute)

	if err != nil {
		log.Fatalf("Error init cache: %v", err)
	}
	var wg sync.WaitGroup

	cfgBroker := config.LoadConfigKafka()
	var brokerCfg []string
	brokerCfg = append(brokerCfg, cfgBroker.Broker)
	broker := kafka.NewKafkaBroker(brokerCfg, cfgBroker.GroupID)
	broker.CreateTopic("orders", 1, 1)

	svc := service.NewOrderService(orderCache, storage)
	handler := handler.NewHandler(svc)
	router := router.NewRouter(handler)

	srv := http.Server{
		Addr:    config.LoadConfigAddr(),
		Handler: router,
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := broker.ConsumeOrder(ctx, "orders", svc.HandleKafkaMessage); err != nil {
			log.Fatalf("kafka error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Listen and running :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Server will shutdown gracefully...")
	err = svc.Shutdown(ctx)
	if err != nil {
		log.Printf("Error shutting down service: %v", err)
	}

	wg.Wait()
}
