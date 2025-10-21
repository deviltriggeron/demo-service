package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	e "demo-service/internal/entity"
	"demo-service/internal/service"
)

type OrderHandler struct {
	svc service.OrderService
}

func NewHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (s *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	order, err := s.svc.Get(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := s.svc.GetAll()
	if err != nil {
		http.Error(w, "Error fetching orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *OrderHandler) Insert(w http.ResponseWriter, r *http.Request) {
	var order e.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Printf("decode error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := s.svc.InsertOrder(order); err != nil {
		http.Error(w, "Error saving order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order e.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.svc.UpdateOrder(order); err != nil {
		http.Error(w, "update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"new_uid": order.OrderUID,
		"message": "order updated successfully",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *OrderHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
