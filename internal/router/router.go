package router

import (
	"github.com/gorilla/mux"

	"demo-service/internal/handler"
	"demo-service/internal/middleware"
)

func NewRouter(h *handler.OrderHandler) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/order/{id}", h.GetOrder).Methods("GET")

	r.HandleFunc("/orders", h.GetOrders).Methods("GET")

	r.HandleFunc("/", h.Index).Methods("GET")

	r.HandleFunc("/insert", h.Insert).Methods("POST")

	r.HandleFunc("/update", h.Update).Methods("PUT")

	return r
}
