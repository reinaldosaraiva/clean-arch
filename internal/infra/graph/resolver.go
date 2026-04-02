package graph

import "github.com/reinaldosaraiva/clean-arch/internal/usecase"

type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}
