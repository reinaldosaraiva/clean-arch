<div align="center">

# clean-arch

[![Go Version](https://img.shields.io/badge/go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/reinaldosaraiva/clean-arch)](https://goreportcard.com/report/github.com/reinaldosaraiva/clean-arch)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![REST](https://img.shields.io/badge/REST-8000-brightgreen)](#rest)
[![gRPC](https://img.shields.io/badge/gRPC-50051-blue)](#grpc)
[![GraphQL](https://img.shields.io/badge/GraphQL-8080-e10098?logo=graphql&logoColor=white)](#graphql)

**Go implementation of the Full Cycle Clean Architecture challenge.**

A single `ListOrders` use case served simultaneously over REST, gRPC, and GraphQL — the domain layer knows none of them.

</div>

---

## Overview

clean-arch demonstrates strict layer separation using the Clean Architecture pattern. The domain defines a repository interface. The use cases depend on that interface. Three independent delivery mechanisms — REST, gRPC, and GraphQL — each call the same use case without any shared infrastructure code.

The project ships as a single `docker compose up --build` command: MySQL starts, migrations run, and all three servers come up on their respective ports.

---

## Getting started

Prerequisites: Docker and Docker Compose.

```bash
git clone https://github.com/reinaldosaraiva/clean-arch.git
cd clean-arch
docker compose up --build
```

That is all. The entrypoint script handles MySQL readiness, runs the migration, and starts the server.

---

## Interfaces

### REST

Port **8000** — Chi router, standard `net/http`.

```bash
# Create an order
curl -s -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"ID":"o1","Price":100.50,"Tax":10.50}' | jq

# List all orders
curl -s http://localhost:8000/order | jq
```

### GraphQL

Port **8080** — gqlgen, schema-first.

Open **http://localhost:8080** for the interactive playground.

```bash
# Create
curl -s -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation{createOrder(input:{id:\"o2\",price:200,tax:20}){id finalPrice}}"}' | jq

# List
curl -s -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{listOrders{id price tax finalPrice}}"}' | jq
```

### gRPC

Port **50051** — protobuf, server reflection enabled (no `.proto` file needed locally).

Requires [`grpcurl`](https://github.com/fullstorydev/grpcurl):

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

The dependency rule is enforced at every layer: outer rings import inner rings, never the reverse.

```
┌──────────────────────────────────────────────────────────┐
│  Delivery (infra)                                        │
│  ┌─────────┐   ┌──────────┐   ┌───────────────────────┐ │
│  │  REST   │   │   gRPC   │   │       GraphQL         │ │
│  └────┬────┘   └────┬─────┘   └──────────┬────────────┘ │
│       │             │                    │               │
│       └─────────────┴────────────────────┘               │
│                          │                               │
│                          ▼                               │
│  ┌───────────────────────────────────────────────────┐   │
│  │  Use Cases                                        │   │
│  │  CreateOrderUseCase · ListOrdersUseCase           │   │
│  └────────────────────────┬──────────────────────────┘   │
│                           │                              │
│                           ▼                              │
│  ┌───────────────────────────────────────────────────┐   │
│  │  Domain                                           │   │
│  │  Order entity · OrderRepositoryInterface          │   │
│  └───────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────┘
                           │
                           ▼ (implements)
               ┌───────────────────────┐
               │  MySQL Repository     │
               │  (infra/database)     │
               └───────────────────────┘
```

### Directory structure

```
cmd/ordersystem/        Wiring: DB, RabbitMQ, use cases, three servers
configs/                Viper-based environment config
internal/
  entity/               Domain: Order, Validate, CalculateFinalPrice
  usecase/              Application: CreateOrder, ListOrders
  infra/
    database/           MySQL: Save, GetTotal, GetAll
    web/                REST handlers (Chi)
    grpc/               gRPC service + protobuf definitions
    graph/              GraphQL resolvers + gqlgen schema
  event/                Domain event: OrderCreated
pkg/events/             Async event dispatcher (goroutines + WaitGroup)
migrations/             SQL: CREATE TABLE orders
scripts/                entrypoint.sh: wait → migrate → exec
```

### Key design decisions

- **OrderRepositoryInterface** is defined in `internal/entity/` — the domain owns the contract, not the infra.
- **CreateOrderUseCase** instantiates a new `OrderCreated` event per call to avoid shared-state race conditions under concurrency.
- **gRPC proto fields** use `double` instead of `float` to preserve financial precision across wire serialization.
- **entrypoint.sh** uses `mysqladmin ping` in a loop so the app never starts before MySQL is ready, removing the need for `depends_on: condition: service_healthy` hacks.

---

## Configuration

All values are provided via environment variables. `docker-compose.yaml` sets them all — no `.env` file is required to run the project.

| Variable | Default | |
|---|---|---|
| `DB_DRIVER` | `mysql` | Database driver |
| `DB_HOST` | `mysql` | MySQL hostname |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASSWORD` | `root` | MySQL password |
| `DB_NAME` | `orders` | Database name |
| `WEB_SERVER_PORT` | `:8000` | REST listen address |
| `GRPC_SERVER_PORT` | `50051` | gRPC listen port |
| `GRAPHQL_SERVER_PORT` | `8080` | GraphQL listen port |
| `RABBITMQ_DSN` | `amqp://guest:guest@rabbitmq:5672/` | RabbitMQ connection string |

---

## Stack

| Component | Technology |
|---|---|
| Language | Go 1.25 |
| Database | MySQL 8.0 |
| Messaging | RabbitMQ 3 — [`amqp091-go`](https://github.com/rabbitmq/amqp091-go) |
| HTTP router | [Chi](https://github.com/go-chi/chi) |
| gRPC | [`google.golang.org/grpc`](https://pkg.go.dev/google.golang.org/grpc) + protobuf |
| GraphQL | [gqlgen](https://github.com/99designs/gqlgen) (schema-first) |
| Config | [Viper](https://github.com/spf13/viper) |
| Container | Docker multi-stage + Compose |

---

## License

MIT
