package service

import (
	"context"

	"github.com/reinaldosaraiva/clean-arch/internal/infra/grpc/pb"
	"github.com/reinaldosaraiva/clean-arch/internal/usecase"
)

type OrderGrpcService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewOrderGrpcService(
	createOrderUseCase usecase.CreateOrderUseCase,
	listOrdersUseCase usecase.ListOrdersUseCase,
) *OrderGrpcService {
	return &OrderGrpcService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderGrpcService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: in.Price,
		Tax:   in.Tax,
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

func (s *OrderGrpcService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	output, err := s.ListOrdersUseCase.Execute(usecase.ListOrdersInputDTO{})
	if err != nil {
		return nil, err
	}
	var orders []*pb.CreateOrderResponse
	for _, o := range output {
		orders = append(orders, &pb.CreateOrderResponse{
			Id:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		})
	}
	return &pb.ListOrdersResponse{Orders: orders}, nil
}
