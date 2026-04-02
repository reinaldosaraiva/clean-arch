package event

import "time"

type OrderCreated struct {
	Name    string
	Payload interface{}
}

func NewOrderCreated() *OrderCreated {
	return &OrderCreated{Name: "OrderCreated"}
}

func (o *OrderCreated) GetName() string          { return o.Name }
func (o *OrderCreated) GetDateTime() time.Time   { return time.Now() }
func (o *OrderCreated) GetPayload() interface{}  { return o.Payload }
func (o *OrderCreated) SetPayload(p interface{}) { o.Payload = p }
