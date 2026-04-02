package usecase

import (
	"github.com/reinaldosaraiva/clean-arch/internal/entity"
	"github.com/reinaldosaraiva/clean-arch/internal/event"
	"github.com/reinaldosaraiva/clean-arch/pkg/events"
)

type OrderInputDTO struct {
	ID    string
	Price float64
	Tax   float64
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository  entity.OrderRepositoryInterface
	EventDispatcher  events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	orderRepository entity.OrderRepositoryInterface,
	eventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
		EventDispatcher: eventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return OrderOutputDTO{}, err
	}
	if err := order.CalculateFinalPrice(); err != nil {
		return OrderOutputDTO{}, err
	}
	if err := c.OrderRepository.Save(order); err != nil {
		return OrderOutputDTO{}, err
	}
	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}
	// Instantiate a new event per execution to avoid shared-state race condition
	orderCreatedEvent := event.NewOrderCreated()
	orderCreatedEvent.SetPayload(dto)
	c.EventDispatcher.Dispatch(orderCreatedEvent)
	return dto, nil
}
