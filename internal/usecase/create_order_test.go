package usecase

import (
	"testing"
	"time"
	"github.com/reinaldosaraiva/clean-arch/pkg/events"
)

type MockEvent struct {
	name    string
	payload any
}

func (e *MockEvent) GetName() string        { return e.name }
func (e *MockEvent) GetDateTime() time.Time { return time.Now() }
func (e *MockEvent) GetPayload() any        { return e.payload }
func (e *MockEvent) SetPayload(p any)       { e.payload = p }

type MockEventDispatcher struct{}

func (d *MockEventDispatcher) Register(eventName string, handler events.EventHandlerInterface) error {
	return nil
}
func (d *MockEventDispatcher) Dispatch(event events.EventInterface) error { return nil }
func (d *MockEventDispatcher) Remove(eventName string, handler events.EventHandlerInterface) error {
	return nil
}
func (d *MockEventDispatcher) Has(eventName string, handler events.EventHandlerInterface) bool {
	return false
}
func (d *MockEventDispatcher) Clear() error { return nil }

func TestCreateOrderUseCase_Execute(t *testing.T) {
	repo := &MockOrderRepository{}
	event := &MockEvent{name: "OrderCreated"}
	dispatcher := &MockEventDispatcher{}
	uc := NewCreateOrderUseCase(repo, event, dispatcher)

	input := OrderInputDTO{ID: "order-1", Price: 100.0, Tax: 10.0}
	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID != "order-1" {
		t.Errorf("expected ID order-1, got %s", output.ID)
	}
	if output.FinalPrice != 110.0 {
		t.Errorf("expected FinalPrice 110.0, got %f", output.FinalPrice)
	}
}
