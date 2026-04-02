package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/reinaldosaraiva/clean-arch/internal/usecase"
)

const maxBodyBytes = 1 << 20 // 1 MB

type WebOrderHandler struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
}

func NewWebOrderHandler(createOrderUseCase usecase.CreateOrderUseCase) *WebOrderHandler {
	return &WebOrderHandler{CreateOrderUseCase: createOrderUseCase}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var dto usecase.OrderInputDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	output, err := h.CreateOrderUseCase.Execute(dto)
	if err != nil {
		log.Printf("create order error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		log.Printf("encode response error: %v", err)
	}
}

type WebListOrderHandler struct {
	ListOrdersUseCase usecase.ListOrdersUseCase
}

func NewWebListOrderHandler(listOrdersUseCase usecase.ListOrdersUseCase) *WebListOrderHandler {
	return &WebListOrderHandler{ListOrdersUseCase: listOrdersUseCase}
}

func (h *WebListOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListOrdersUseCase.Execute(usecase.ListOrdersInputDTO{})
	if err != nil {
		log.Printf("list orders error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if output == nil {
		output = []usecase.OrderOutputDTO{}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(output); err != nil {
		log.Printf("encode response error: %v", err)
	}
}
