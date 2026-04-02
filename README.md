# Clean Architecture — Go

Implementacao de Clean Architecture em Go com o use case **ListOrders** exposto simultaneamente via **REST**, **gRPC** e **GraphQL**.

## Execucao

```bash
docker compose up --build
```

Esse unico comando sobe o banco de dados, executa as migrations e inicia a aplicacao.

## Portas

| Protocolo | Porta | Endpoint |
|-----------|-------|----------|
| REST      | 8000  | `GET /order`, `POST /order` |
| gRPC      | 50051 | `pb.OrderService/ListOrders`, `pb.OrderService/CreateOrder` |
| GraphQL   | 8080  | `POST /query` (API), `GET /` (Playground) |

## Como testar

### REST

```bash
# Criar order
curl -s -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"ID":"order-1","Price":100.5,"Tax":10.5}' | jq

# Listar orders
curl -s http://localhost:8000/order | jq
```

### GraphQL

Acesse o Playground interativo: http://localhost:8080

```graphql
# Criar order
mutation {
  createOrder(input: { id: "order-2", price: 200.0, tax: 20.0 }) {
    id
    price
    tax
    finalPrice
  }
}

# Listar orders
{
  listOrders {
    id
    price
    tax
    finalPrice
  }
}
```

### gRPC

Requer `grpcurl` instalado (`brew install grpcurl`):

```bash
# Criar order
grpcurl -plaintext \
  -d '{"id":"order-3","price":50.0,"tax":5.0}' \
  localhost:50051 pb.OrderService/CreateOrder

# Listar orders
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
```

## Arquitetura

```
cmd/ordersystem/     — entrypoint (main.go com DI manual)
configs/             — Viper config loader
internal/
  entity/            — Order entity, OrderRepositoryInterface
  usecase/           — CreateOrderUseCase, ListOrdersUseCase
  event/             — OrderCreated domain event, OrderCreatedHandler
  infra/
    database/        — MySQL OrderRepository
    web/             — REST handlers (Chi)
    grpc/            — proto, pb gerado, OrderGrpcService
    graph/           — GraphQL schema, resolvers (gqlgen)
pkg/events/          — EventDispatcher
migrations/          — SQL migrations
scripts/             — entrypoint.sh (wait-for-mysql + migrate + start)
```

## Dependencias

- Go 1.22+
- Docker e Docker Compose
