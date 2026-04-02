# Clean Architecture — Go

Implementação do desafio **Full Cycle — Clean Architecture** em Go.

O objetivo é demonstrar o **desacoplamento** da arquitetura: um único use case (`ListOrders`) é exposto simultaneamente por três interfaces de comunicação independentes — REST, gRPC e GraphQL — sem que a camada de domínio conheça nenhuma delas.

---

## Execução

> Pré-requisito: Docker e Docker Compose instalados.

```bash
docker compose up --build
```

Esse único comando:
1. Sobe o MySQL 8.0 e aguarda o healthcheck passar
2. Executa as migrations automaticamente (`CREATE TABLE orders`)
3. Inicia a aplicação Go nas três portas simultaneamente

Nenhum outro comando é necessário.

---

## Serviços e portas

| Protocolo | Porta | Endpoints |
|-----------|-------|-----------|
| REST      | 8000  | `POST /order` — criar order |
|           |       | `GET /order` — listar orders |
| gRPC      | 50051 | `pb.OrderService/CreateOrder` |
|           |       | `pb.OrderService/ListOrders` |
| GraphQL   | 8080  | `POST /query` — API |
|           |       | `GET /` — Playground interativo |

---

## Testando cada protocolo

### REST

```bash
# Criar uma order
curl -s -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"ID":"order-1","Price":100.50,"Tax":10.50}' | jq

# Listar todas as orders
curl -s http://localhost:8000/order | jq
```

### GraphQL

Acesse o Playground em **http://localhost:8080** e execute:

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
query {
  listOrders {
    id
    price
    tax
    finalPrice
  }
}
```

Ou via HTTP direto:

```bash
curl -s -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ listOrders { id price tax finalPrice } }"}' | jq
```

### gRPC

Requer [`grpcurl`](https://github.com/fullstorydev/grpcurl) (`brew install grpcurl`):

```bash
# Criar order
grpcurl -plaintext \
  -d '{"id":"order-3","price":50.0,"tax":5.0}' \
  localhost:50051 pb.OrderService/CreateOrder

# Listar orders
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
```

> O servidor gRPC tem reflection habilitada, então `grpcurl` funciona sem precisar do arquivo `.proto`.

---

## Arquitetura

O projeto segue **Clean Architecture** com separação estrita entre camadas. A regra central é que o domínio nunca depende de infraestrutura.

```
cmd/ordersystem/
  main.go              # Wiring manual: MySQL, RabbitMQ, UseCases, 3 servidores

configs/
  config.go            # Carrega variáveis de ambiente via Viper

internal/
  entity/
    order.go           # Entidade Order com validação e CalculateFinalPrice
    interface.go       # OrderRepositoryInterface (Save, GetTotal, GetAll)

  usecase/
    create_order.go    # CreateOrderUseCase — salva e emite evento OrderCreated
    list_orders.go     # ListOrdersUseCase — delega ao repositório

  event/
    order_created.go              # Evento de domínio OrderCreated
    handler/order_created_handler.go  # Publica no RabbitMQ

  infra/
    database/
      order_repository.go  # MySQL: Save, GetTotal, GetAll

    web/
      order_handler.go     # Handlers REST (POST e GET /order)
      webserver/           # Chi router

    grpc/
      protofiles/order.proto   # Definição do serviço gRPC
      pb/                      # Código gerado pelo protoc
      service/order_service.go # Implementação dos RPCs

    graph/
      schema.graphqls      # Schema GraphQL
      schema.resolvers.go  # Resolvers (createOrder, listOrders)
      generated.go         # Gerado pelo gqlgen

pkg/events/
  event_dispatcher.go    # Dispatcher assíncrono com goroutines

migrations/
  001_create_orders.sql  # CREATE TABLE orders

scripts/
  entrypoint.sh          # Aguarda MySQL, roda migrations, inicia app
```

### Fluxo de uma requisição ListOrders

```
Cliente (REST / gRPC / GraphQL)
        │
        ▼
   Handler (infra)          ← conhece o UseCase, não o domínio
        │
        ▼
  ListOrdersUseCase         ← conhece a interface do repositório
        │
        ▼
  OrderRepositoryInterface  ← definida no domínio
        │
        ▼
  OrderRepository (MySQL)   ← implementação concreta na infra
```

---

## Variáveis de ambiente

Todas as configurações são injetadas via variáveis de ambiente (veja `.env.example`). Dentro do Docker, o `docker-compose.yaml` já define todos os valores necessários.

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `DB_DRIVER` | `mysql` | Driver do banco |
| `DB_HOST` | `mysql` | Host do MySQL |
| `DB_PORT` | `3306` | Porta do MySQL |
| `DB_USER` | `root` | Usuário |
| `DB_PASSWORD` | `root` | Senha |
| `DB_NAME` | `orders` | Nome do banco |
| `WEB_SERVER_PORT` | `:8000` | Porta REST |
| `GRPC_SERVER_PORT` | `50051` | Porta gRPC |
| `GRAPHQL_SERVER_PORT` | `8080` | Porta GraphQL |
| `RABBITMQ_DSN` | `amqp://guest:guest@rabbitmq:5672/` | DSN do RabbitMQ |

---

## Tecnologias

| Tecnologia | Uso |
|------------|-----|
| Go 1.25 | Linguagem |
| MySQL 8.0 | Banco de dados |
| RabbitMQ 3 | Mensageria (evento OrderCreated) |
| Chi | Router HTTP |
| gRPC / protobuf | Comunicação gRPC |
| gqlgen | Geração de código GraphQL |
| Viper | Configuração via env |
| Docker / Compose | Infraestrutura |
