package usecase

import "github.com/reinaldosaraiva/clean-arch/internal/entity"

type ListOrdersInputDTO struct{}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(repo entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{OrderRepository: repo}
}

func (u *ListOrdersUseCase) Execute(input ListOrdersInputDTO) ([]OrderOutputDTO, error) {
	orders, err := u.OrderRepository.GetAll()
	if err != nil {
		return nil, err
	}
	var out []OrderOutputDTO
	for _, o := range orders {
		out = append(out, OrderOutputDTO{
			ID:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		})
	}
	return out, nil
}
