# clean-arch

[![Go Version](https://img.shields.io/badge/go-1.25-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue)](#license)
[![Build](https://img.shields.io/badge/build-docker%20compose-2496ED?logo=docker)](https://docs.docker.com/compose/)

A **Full Cycle Clean Architecture** challenge in Go — one `ListOrders` use case delivered simultaneously over **REST**, **gRPC**, and **GraphQL**, without the domain layer knowing any of them.

---

## Quick start

```bash
git clone https://github.com/reinaldosaraiva/clean-arch.git
cd clean-arch
docker compose up --build
```

That single command:
- starts MySQL 8.0 and waits for it to be healthy
- runs the database migration (`CREATE TABLE orders`)
- starts the Go server on all three interfaces

---

## Endpoints

| Interface | Port  | Operation       | Address |
|-----------|-------|-----------------|---------|
| REST      | 8000  | Create order    | `POST /order` |
|           |       | List orders     | `GET /order` |
| gRPC      | 50051 | Create order    | `pb.OrderService/CreateOrder` |
|           |       | List orders     | `pb.OrderService/ListOrders` |
| GraphQL   | 8080  | Playground      | `GET /` |
|           |       | API             | `POST /query` |

---

## Usage

### REST

```bash
# Create
curl -s -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"ID":"o1","Price":100.50,"Tax":10.50}' | jq

# List
curl -s http://localhost:8000/order | jq
```

### GraphQL

Open **http://localhost:8080** for the interactive playground, or call the API directly:

```bash
# Create
curl -s -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation { createOrder(input:{id:\"o2\",price:200,tax:20}){ id finalPrice } }"}' | jq

# List
curl -s -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ listOrders { id price tax finalPrice } }"}' | jq
```

### gRPC

Requires [`grpcurl`](https://github.com/fullstorydev/grpcurl) — reflection is enabled, no `.proto` file needed:

```bash
# Create
grpcurl -plaintext \
  -d '{"id":"o3","price":50,"tax":5}' \
  localhost:50051 pb.OrderService/CreateOrder

# List
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
```

---

## Architecture

The project follows Clean Architecture with a strict unidirectional dependency rule: outer layers depend on inner ones, never the reverse.

```
internal/
├── entity/          # Domain: Order, validation, CalculateFinalPrice
├── usecase/         # Application: CreateOrder, ListOrders
├── infra/
│   ├── database/    # MySQL repository (implements entity.OrderRepositoryInterface)
│   ├── web/         # REST handlers (Chi)
│   ├── grpc/        # gRPC service (protobuf)
│   └── graph/       # GraphQL resolvers (gqlgen)
└── event/           # Domain event: OrderCreated

cmd/ordersystem/     # Wiring: DB, RabbitMQ, use cases, three servers
pkg/events/          # Async event dispatcher
migrations/          # SQL migrations
```

**Dependency direction:**

```
infra  →  usecase  →  entity
              ↑
         (interfaces only)
```

The domain (`entity/`) defines `OrderRepositoryInterface`. The use cases depend on that interface. The MySQL repository implements it. No layer imports anything above itself.

---

## Configuration

All values are injected via environment variables. Inside Docker, `docker-compose.yaml` provides them all — no `.env` file required.

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_DRIVER` | `mysql` | Database driver |
| `DB_HOST` | `mysql` | MySQL hostname |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASSWORD` | `root` | MySQL password |
| `DB_NAME` | `orders` | Database name |
| `WEB_SERVER_PORT` | `:8000` | REST port |
| `GRPC_SERVER_PORT` | `50051` | gRPC port |
| `GRAPHQL_SERVER_PORT` | `8080` | GraphQL port |
| `RABBITMQ_DSN` | `amqp://guest:guest@rabbitmq:5672/` | RabbitMQ DSN |

---

## Stack

| | |
|---|---|
| Language | Go 1.25 |
| Database | MySQL 8.0 |
| Messaging | RabbitMQ 3 (`rabbitmq/amqp091-go`) |
| HTTP router | [Chi](https://github.com/go-chi/chi) |
| gRPC | `google.golang.org/grpc` + protobuf |
| GraphQL | [gqlgen](https://github.com/99designs/gqlgen) |
| Config | [Viper](https://github.com/spf13/viper) |
| Container | Docker + Compose |

---

## License

MIT
