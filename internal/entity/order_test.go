package entity

import (
	"testing"
)

func TestNewOrder_ValidData(t *testing.T) {
	order, err := NewOrder("123", 10.0, 1.5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.ID != "123" {
		t.Errorf("expected ID '123', got '%s'", order.ID)
	}
	if order.Price != 10.0 {
		t.Errorf("expected Price 10.0, got %f", order.Price)
	}
	if order.Tax != 1.5 {
		t.Errorf("expected Tax 1.5, got %f", order.Tax)
	}
}

func TestNewOrder_EmptyID(t *testing.T) {
	_, err := NewOrder("", 10.0, 1.5)
	if err == nil {
		t.Fatal("expected error for empty ID, got nil")
	}
	if err.Error() != "invalid id" {
		t.Errorf("expected 'invalid id', got '%s'", err.Error())
	}
}

func TestNewOrder_InvalidPrice(t *testing.T) {
	_, err := NewOrder("123", 0, 1.5)
	if err == nil {
		t.Fatal("expected error for price <= 0, got nil")
	}
	if err.Error() != "invalid price" {
		t.Errorf("expected 'invalid price', got '%s'", err.Error())
	}

	_, err = NewOrder("123", -5.0, 1.5)
	if err == nil {
		t.Fatal("expected error for negative price, got nil")
	}
	if err.Error() != "invalid price" {
		t.Errorf("expected 'invalid price', got '%s'", err.Error())
	}
}

func TestNewOrder_InvalidTax(t *testing.T) {
	_, err := NewOrder("123", 10.0, 0)
	if err == nil {
		t.Fatal("expected error for tax <= 0, got nil")
	}
	if err.Error() != "invalid tax" {
		t.Errorf("expected 'invalid tax', got '%s'", err.Error())
	}

	_, err = NewOrder("123", 10.0, -1.0)
	if err == nil {
		t.Fatal("expected error for negative tax, got nil")
	}
	if err.Error() != "invalid tax" {
		t.Errorf("expected 'invalid tax', got '%s'", err.Error())
	}
}

func TestCalculateFinalPrice(t *testing.T) {
	order, err := NewOrder("123", 10.0, 1.5)
	if err != nil {
		t.Fatalf("expected no error creating order, got %v", err)
	}
	if err := order.CalculateFinalPrice(); err != nil {
		t.Fatalf("expected no error calculating final price, got %v", err)
	}
	expected := 11.5
	if order.FinalPrice != expected {
		t.Errorf("expected FinalPrice %f, got %f", expected, order.FinalPrice)
	}
}
