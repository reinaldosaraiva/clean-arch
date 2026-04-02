package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/reinaldosaraiva/clean-arch/internal/entity"
	"github.com/reinaldosaraiva/clean-arch/internal/usecase"
	"github.com/reinaldosaraiva/clean-arch/pkg/events"
)

// MockOrderRepository is an in-memory order repository for tests.
type MockOrderRepository struct {
	orders []entity.Order
}

func (m *MockOrderRepository) Save(o *entity.Order) error {
	m.orders = append(m.orders, *o)
	return nil
}

func (m *MockOrderRepository) GetTotal() (int, error) {
	return len(m.orders), nil
}

func (m *MockOrderRepository) GetAll() ([]entity.Order, error) {
	return m.orders, nil
}

// MockEvent implements events.EventInterface.
type MockEvent struct {
	payload any
}

func (e *MockEvent) GetName() string        { return "OrderCreated" }
func (e *MockEvent) GetDateTime() time.Time { return time.Now() }
func (e *MockEvent) GetPayload() any        { return e.payload }
func (e *MockEvent) SetPayload(p any)       { e.payload = p }

// MockDispatcher implements events.EventDispatcherInterface.
type MockDispatcher struct{}

func (d *MockDispatcher) Register(n string, h events.EventHandlerInterface) error { return nil }
func (d *MockDispatcher) Dispatch(e events.EventInterface) error                  { return nil }
func (d *MockDispatcher) Remove(n string, h events.EventHandlerInterface) error   { return nil }
func (d *MockDispatcher) Has(n string, h events.EventHandlerInterface) bool       { return false }
func (d *MockDispatcher) Clear() error                                             { return nil }

// Ensure MockDispatcher satisfies the interface at compile time.
var _ events.EventDispatcherInterface = (*MockDispatcher)(nil)

// Silence unused import of sync.
var _ = sync.WaitGroup{}

func TestWebOrderHandler_Create(t *testing.T) {
	repo := &MockOrderRepository{}
	uc := *usecase.NewCreateOrderUseCase(repo, &MockDispatcher{})
	handler := NewWebOrderHandler(uc)

	body := `{"ID":"order-1","Price":100.0,"Tax":10.0}`
	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	var out usecase.OrderOutputDTO
	json.NewDecoder(w.Body).Decode(&out)
	if out.FinalPrice != 110.0 {
		t.Errorf("expected FinalPrice 110.0, got %f", out.FinalPrice)
	}
}

func TestWebListOrderHandler_List(t *testing.T) {
	repo := &MockOrderRepository{
		orders: []entity.Order{
			{ID: "1", Price: 10, Tax: 1, FinalPrice: 11},
		},
	}
	uc := *usecase.NewListOrdersUseCase(repo)
	handler := NewWebListOrderHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/order", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var out []usecase.OrderOutputDTO
	json.NewDecoder(w.Body).Decode(&out)
	if len(out) != 1 {
		t.Errorf("expected 1 order, got %d", len(out))
	}
}

func TestWebListOrderHandler_ListEmpty(t *testing.T) {
	repo := &MockOrderRepository{}
	uc := *usecase.NewListOrdersUseCase(repo)
	handler := NewWebListOrderHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/order", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var out []usecase.OrderOutputDTO
	json.NewDecoder(w.Body).Decode(&out)
	if len(out) != 0 {
		t.Errorf("expected 0 orders, got %d", len(out))
	}
}
