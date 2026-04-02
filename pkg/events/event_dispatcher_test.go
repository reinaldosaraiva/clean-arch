package events

import (
	"sync"
	"testing"
	"time"
)

// TestEvent is a mock event for testing
type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string        { return e.Name }
func (e *TestEvent) GetDateTime() time.Time { return time.Now() }
func (e *TestEvent) GetPayload() interface{} { return e.Payload }
func (e *TestEvent) SetPayload(p interface{}) { e.Payload = p }

// TestEventHandler is a mock handler for testing
type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
}

func TestEventDispatcher_Register(t *testing.T) {
	ed := NewEventDispatcher()
	h1 := &TestEventHandler{ID: 1}

	err := ed.Register("order.created", h1)
	if err != nil {
		t.Fatalf("expected no error on first registration, got %v", err)
	}

	if !ed.Has("order.created", h1) {
		t.Error("expected handler to be registered")
	}
}

func TestEventDispatcher_Register_DuplicateHandler(t *testing.T) {
	ed := NewEventDispatcher()
	h1 := &TestEventHandler{ID: 1}

	_ = ed.Register("order.created", h1)
	err := ed.Register("order.created", h1)
	if err == nil {
		t.Fatal("expected error on duplicate registration, got nil")
	}
	if err != ErrHandlerAlreadyRegistered {
		t.Errorf("expected ErrHandlerAlreadyRegistered, got %v", err)
	}
}

func TestEventDispatcher_Dispatch(t *testing.T) {
	ed := NewEventDispatcher()
	h1 := &TestEventHandler{ID: 1}
	h2 := &TestEventHandler{ID: 2}

	_ = ed.Register("order.created", h1)
	_ = ed.Register("order.created", h2)

	event := &TestEvent{Name: "order.created", Payload: "test-payload"}
	err := ed.Dispatch(event)
	if err != nil {
		t.Fatalf("expected no error dispatching event, got %v", err)
	}
}

func TestEventDispatcher_Remove(t *testing.T) {
	ed := NewEventDispatcher()
	h1 := &TestEventHandler{ID: 1}

	_ = ed.Register("order.created", h1)
	if !ed.Has("order.created", h1) {
		t.Fatal("handler should be registered before removal")
	}

	_ = ed.Remove("order.created", h1)
	if ed.Has("order.created", h1) {
		t.Error("handler should have been removed")
	}
}

func TestEventDispatcher_Has(t *testing.T) {
	ed := NewEventDispatcher()
	h1 := &TestEventHandler{ID: 1}
	h2 := &TestEventHandler{ID: 2}

	_ = ed.Register("order.created", h1)

	if !ed.Has("order.created", h1) {
		t.Error("expected Has to return true for registered handler")
	}
	if ed.Has("order.created", h2) {
		t.Error("expected Has to return false for unregistered handler")
	}
	if ed.Has("other.event", h1) {
		t.Error("expected Has to return false for unknown event")
	}
}

func TestEventDispatcher_Clear(t *testing.T) {
	ed := NewEventDispatcher()
	h1 := &TestEventHandler{ID: 1}
	h2 := &TestEventHandler{ID: 2}

	_ = ed.Register("order.created", h1)
	_ = ed.Register("order.updated", h2)

	err := ed.Clear()
	if err != nil {
		t.Fatalf("expected no error on clear, got %v", err)
	}

	if ed.Has("order.created", h1) {
		t.Error("expected handler to be cleared")
	}
	if ed.Has("order.updated", h2) {
		t.Error("expected handler to be cleared")
	}
}

func TestEventDispatcher_MultipleHandlers(t *testing.T) {
	ed := NewEventDispatcher()
	h1 := &TestEventHandler{ID: 1}
	h2 := &TestEventHandler{ID: 2}

	_ = ed.Register("order.created", h1)
	_ = ed.Register("order.created", h2)

	if !ed.Has("order.created", h1) {
		t.Error("expected h1 to be registered")
	}
	if !ed.Has("order.created", h2) {
		t.Error("expected h2 to be registered")
	}
}
