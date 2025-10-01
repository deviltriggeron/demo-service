package server

import (
	model "demo-service/internal/model"
	service "demo-service/internal/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	svc *service.OrderService
}

func NewServer(svc *service.OrderService) *Server {
	return &Server{svc: svc}
}

func (s *Server) Run() error {
	r := mux.NewRouter()

	r.HandleFunc("/order/{id}", s.getOrder).Methods("GET")

	r.HandleFunc("/orders", s.getOrders).Methods("GET")

	r.HandleFunc("/", s.index).Methods("GET")

	r.HandleFunc("/insert", s.insert).Methods("POST")

	r.HandleFunc("/update", s.update).Methods("PUT")

	log.Println("HTTP server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	order, err := s.svc.Get(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (s *Server) getOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := s.svc.GetAll()
	if err != nil {
		http.Error(w, "Error fetching orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (s *Server) insert(w http.ResponseWriter, r *http.Request) {
	var order model.Order

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
	json.NewEncoder(w).Encode(order)
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order model.Order
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
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"new_uid": order.OrderUID,
		"message": "order updated successfully",
	})
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
