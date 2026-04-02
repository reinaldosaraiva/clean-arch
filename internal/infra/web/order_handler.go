package web

import (
	"encoding/json"
	"net/http"

	"github.com/reinaldosaraiva/clean-arch/internal/usecase"
)

type WebOrderHandler struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
}

func NewWebOrderHandler(createOrderUseCase usecase.CreateOrderUseCase) *WebOrderHandler {
	return &WebOrderHandler{CreateOrderUseCase: createOrderUseCase}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.CreateOrderUseCase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if output == nil {
		output = []usecase.OrderOutputDTO{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
