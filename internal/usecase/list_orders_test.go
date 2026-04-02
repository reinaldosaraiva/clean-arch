package usecase

import (
	"testing"
	"github.com/reinaldosaraiva/clean-arch/internal/entity"
)

type MockOrderRepository struct {
	orders []entity.Order
}

func (m *MockOrderRepository) Save(order *entity.Order) error {
	m.orders = append(m.orders, *order)
	return nil
}

func (m *MockOrderRepository) GetTotal() (int, error) {
	return len(m.orders), nil
}

func (m *MockOrderRepository) GetAll() ([]entity.Order, error) {
	return m.orders, nil
}

func TestListOrdersUseCase_Empty(t *testing.T) {
	repo := &MockOrderRepository{}
	uc := NewListOrdersUseCase(repo)
	result, err := uc.Execute(ListOrdersInputDTO{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 orders, got %d", len(result))
	}
}

func TestListOrdersUseCase_WithOrders(t *testing.T) {
	repo := &MockOrderRepository{
		orders: []entity.Order{
			{ID: "1", Price: 10.0, Tax: 1.0, FinalPrice: 11.0},
			{ID: "2", Price: 20.0, Tax: 2.0, FinalPrice: 22.0},
		},
	}
	uc := NewListOrdersUseCase(repo)
	result, err := uc.Execute(ListOrdersInputDTO{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(result))
	}
	if result[0].ID != "1" || result[0].FinalPrice != 11.0 {
		t.Errorf("unexpected first order: %+v", result[0])
	}
	if result[1].ID != "2" || result[1].FinalPrice != 22.0 {
		t.Errorf("unexpected second order: %+v", result[1])
	}
}
