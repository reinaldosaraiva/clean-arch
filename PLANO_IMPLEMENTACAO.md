# Plano de Implementacao: Clean Architecture Go
## ListOrders via REST + gRPC + GraphQL com Docker Automation

---

## 1. Contexto e Objetivo

Construir uma aplicacao Go seguindo Clean Architecture com o use case `ListOrders` exposto
simultaneamente via REST, gRPC e GraphQL. A infraestrutura Docker deve ser completamente
automatizada: banco de dados sobe, migrations executam e a aplicacao inicia sem nenhum
comando manual alem de `docker compose up`.

**Repositorio base de referencia:** `https://github.com/devfullcycle/goexpert/tree/main/20-CleanArch`

O repositorio `clean-arch` esta atualmente **vazio** (apenas log de sessao anterior).
Toda a estrutura sera criada do zero.

---

## 2. Decisao de Complexidade

| Criterio | Score |
|----------|-------|
| Criar novos arquivos (Dockerfile, migrations, usecase, proto, schema) | +1 |
| Afeta multiplos modulos (entity, usecase, 3x infra, wire, docker) | +1 |
| Banco de dados + migrations automaticas | +1 |
| Mudanca arquitetural (3 protocolos simultaneos + Docker full automation) | +2 |
| Testes unitarios e integration | +1 |
| Integracao multi-protocolo (REST + gRPC + GraphQL) | +1 |

**SCORE:** 7 | **LEVEL:** COMPLEX | **EFFORT:** max | **MODEL:** opus

---

## 3. Stack Tecnica

| Componente | Tecnologia | Justificativa |
|------------|------------|---------------|
| Linguagem | Go 1.22+ | Requisito do desafio |
| Banco de dados | MySQL 8.0 | Compativel com template de referencia |
| REST framework | `net/http` + `chi` | Padrao do template de referencia |
| gRPC | `google.golang.org/grpc` + `protoc` | Padrao Go para gRPC |
| GraphQL | `github.com/99designs/gqlgen` | Schema-first, padrao Go |
| DI | `github.com/google/wire` | Padrao do template de referencia |
| Config | `github.com/spf13/viper` | Padrao do template de referencia |
| Migrations | SQL puro via entrypoint | Sem dependencia extra de ferramenta |
| Eventos | `github.com/streadway/amqp` (RabbitMQ) | Ja presente no template |

**Modulo Go:** `github.com/reinaldosaraiva/clean-arch`

**Portas:**
- REST: `:8000`
- gRPC: `:50051`
- GraphQL: `:8080`

---

## 4. Estrutura de Diretorios Final

```
clean-arch/
├── cmd/
│   └── ordersystem/
│       ├── main.go              # Inicializa e sobe os 3 servidores
│       ├── wire.go              # Providers Wire (DI)
│       └── wire_gen.go          # Gerado pelo wire (nao editar manualmente)
├── configs/
│   └── config.go                # Struct de config + Viper loader
├── internal/
│   ├── entity/
│   │   ├── order.go             # Entidade Order + CalculateFinalPrice()
│   │   ├── order_test.go        # Testes unitarios da entidade
│   │   └── interface.go         # OrderRepositoryInterface (Save + GetAll)
│   ├── usecase/
│   │   ├── create_order.go      # CreateOrderUseCase (existente no template)
│   │   ├── create_order_test.go
│   │   ├── list_orders.go       # ListOrdersUseCase (NOVO)
│   │   └── list_orders_test.go  # (NOVO)
│   └── infra/
│       ├── database/
│       │   └── order_repository.go  # MySQL: Save() + GetAll()
│       ├── web/
│       │   ├── order_handler.go     # POST /order + GET /order
│       │   └── webserver/
│       │       └── webserver.go     # Chi router setup
│       ├── grpc/
│       │   ├── protofiles/
│       │   │   └── order.proto      # CreateOrder + ListOrders RPC
│       │   ├── pb/
│       │   │   ├── order.pb.go      # Gerado pelo protoc
│       │   │   └── order_grpc.pb.go # Gerado pelo protoc
│       │   └── service/
│       │       └── order_service.go # CreateOrder + ListOrders handlers
│       └── graph/
│           ├── schema.graphqls      # Mutation createOrder + Query listOrders
│           ├── resolver.go          # Struct Resolver com usecases injetados
│           ├── schema.resolvers.go  # Implementacoes dos resolvers
│           ├── generated.go         # Gerado pelo gqlgen
│           └── model/
│               └── models_gen.go    # Gerado pelo gqlgen
├── pkg/
│   └── events/
│       ├── interface.go
│       ├── event_dispatcher.go
│       └── event_dispatcher_test.go
├── migrations/
│   └── 001_create_orders.sql    # CREATE DATABASE + CREATE TABLE orders
├── scripts/
│   └── entrypoint.sh            # wait-for-mysql + migrate + start app
├── Dockerfile                   # Multi-stage build
├── docker-compose.yaml          # MySQL + RabbitMQ + App com healthcheck
├── gqlgen.yml                   # Config do gqlgen
├── tools.go                     # Wire + gqlgen como tool dependencies
├── .env                         # Variaveis de ambiente (nao commitar em prod)
├── api.http                     # Requisicoes prontas: create + list (3 protocolos)
└── README.md                    # docker compose up + portas + como testar
```

---

## 5. Orchestracao Multi-Agent

```
ORQUESTRADOR (Sonnet 4.6)
├── [FASE 1] golang-pro        → Scaffolding + go.mod + entidade + interfaces
├── [FASE 2] golang-pro        → UseCase ListOrders + Repository GetAll()
├── [FASE 3A] golang-pro       → REST Handler GET /order            ─┐
├── [FASE 3B] golang-pro       → gRPC proto + service ListOrders    ─┤ PARALELO
├── [FASE 3C] golang-pro       → GraphQL schema + resolver          ─┘
├── [FASE 4] golang-pro        → Wire DI update (apos 3A+3B+3C)
├── [FASE 5] golang-pro        → Docker: Dockerfile + compose + migrations + entrypoint
├── [FASE 6] golang-pro        → api.http + README
├── [FASE 7A] code-reviewer    → Revisao Clean Architecture + Go idioms   ─┐ PARALELO
├── [FASE 7B] security-auditor → SQL injection + credenciais + portas      ─┤
└── [FASE 7C] test-runner      → go test ./... + docker compose up         ─┘
```

---

## 6. Fases Detalhadas

---

### FASE 1 — Scaffolding e Dominio

**Agente:** `golang-pro` | **Modo:** foreground | **Isolation:** worktree

#### 1.1 Inicializar modulo Go

```bash
go mod init github.com/reinaldosaraiva/clean-arch
```

#### 1.2 Instalar dependencias

```bash
go get github.com/go-chi/chi/v5
go get github.com/spf13/viper
go get google.golang.org/grpc
go get google.golang.org/protobuf
go get github.com/99designs/gqlgen
go get github.com/google/wire/cmd/wire
go get github.com/go-sql-driver/mysql
go get github.com/streadway/amqp
```

#### 1.3 Criar entidade Order

**Arquivo:** `internal/entity/order.go`

```go
package entity

import "errors"

type Order struct {
    ID         string
    Price      float64
    Tax        float64
    FinalPrice float64
}

func NewOrder(id string, price float64, tax float64) (*Order, error) {
    order := &Order{
        ID:    id,
        Price: price,
        Tax:   tax,
    }
    if err := order.Validate(); err != nil {
        return nil, err
    }
    return order, nil
}

func (o *Order) CalculateFinalPrice() error {
    if err := o.Validate(); err != nil {
        return err
    }
    o.FinalPrice = o.Price + o.Tax
    return nil
}

func (o *Order) Validate() error {
    if o.ID == "" {
        return errors.New("invalid id")
    }
    if o.Price <= 0 {
        return errors.New("invalid price")
    }
    if o.Tax <= 0 {
        return errors.New("invalid tax")
    }
    return nil
}
```

#### 1.4 Definir interface do repositorio

**Arquivo:** `internal/entity/interface.go`

```go
package entity

type OrderRepositoryInterface interface {
    Save(order *Order) error
    GetTotal() (int, error)
    GetAll() ([]Order, error)  // NOVO — necessario para ListOrders
}
```

**Criterio de aceite da Fase 1:**
- [ ] `go build ./internal/entity/` compila sem erros
- [ ] `go test ./internal/entity/...` passa

---

### FASE 2 — Use Case ListOrders

**Agente:** `golang-pro` | **Modo:** foreground

#### 2.1 Implementar use case

**Arquivo:** `internal/usecase/list_orders.go`

```go
package usecase

import "github.com/reinaldosaraiva/clean-arch/internal/entity"

// ListOrdersInputDTO intencional vazio: listagem sem filtros nesta versao
type ListOrdersInputDTO struct{}

type ListOrdersOutputDTO struct {
    Orders []OrderOutputDTO
}

type ListOrdersUseCase struct {
    OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(repo entity.OrderRepositoryInterface) *ListOrdersUseCase {
    return &ListOrdersUseCase{OrderRepository: repo}
}

func (u *ListOrdersUseCase) Execute(input ListOrdersInputDTO) (ListOrdersOutputDTO, error) {
    orders, err := u.OrderRepository.GetAll()
    if err != nil {
        return ListOrdersOutputDTO{}, err
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
    return ListOrdersOutputDTO{Orders: out}, nil
}
```

#### 2.2 Implementar GetAll() no repositorio MySQL

**Arquivo:** `internal/infra/database/order_repository.go`

Adicionar ao repositorio existente:

```go
func (r *OrderRepository) GetAll() ([]entity.Order, error) {
    rows, err := r.Db.Query("SELECT id, price, tax, final_price FROM orders")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []entity.Order
    for rows.Next() {
        var o entity.Order
        if err := rows.Scan(&o.ID, &o.Price, &o.Tax, &o.FinalPrice); err != nil {
            return nil, err
        }
        orders = append(orders, o)
    }
    return orders, nil
}
```

**Criterio de aceite da Fase 2:**
- [ ] `go build ./internal/usecase/...` compila sem erros
- [ ] `go test ./internal/usecase/...` passa (mock do repositorio)

---

### FASE 3A — Interface REST: GET /order

**Agente:** `golang-pro` | **Modo:** foreground | **Paralelo com 3B e 3C**

#### 3A.1 Adicionar handler no order_handler.go

**Arquivo:** `internal/infra/web/order_handler.go`

```go
// ListOrdersHandler adiciona o novo handler ao WebOrderHandler existente
type WebListOrderHandler struct {
    ListOrdersUseCase usecase.ListOrdersUseCase
}

func NewWebListOrderHandler(useCase usecase.ListOrdersUseCase) *WebListOrderHandler {
    return &WebListOrderHandler{ListOrdersUseCase: useCase}
}

func (h *WebListOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
    output, err := h.ListOrdersUseCase.Execute(usecase.ListOrdersInputDTO{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(output.Orders)
}
```

#### 3A.2 Registrar rota GET /order no webserver

**Arquivo:** `cmd/ordersystem/main.go`

```go
// Adicionar ao setup do chi router:
webserver.AddHandler("GET", "/order", listOrderHandler.Create)
```

**Criterio de aceite da Fase 3A:**
- [ ] `curl -s http://localhost:8000/order` retorna JSON (array, mesmo que vazio)
- [ ] HTTP 200 com `Content-Type: application/json`

---

### FASE 3B — Interface gRPC: ListOrders

**Agente:** `golang-pro` | **Modo:** foreground | **Paralelo com 3A e 3C**

#### 3B.1 Atualizar proto

**Arquivo:** `internal/infra/grpc/protofiles/order.proto`

```proto
syntax = "proto3";

package pb;

option go_package = "github.com/reinaldosaraiva/clean-arch/internal/infra/grpc/pb";

service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);  // NOVO
}

message CreateOrderRequest {
  string id = 1;
  float price = 2;
  float tax = 3;
}

message CreateOrderResponse {
  string id = 1;
  float price = 2;
  float tax = 3;
  float final_price = 4;
}

message ListOrdersRequest {}  // NOVO

message ListOrdersResponse {  // NOVO
  repeated CreateOrderResponse orders = 1;
}
```

#### 3B.2 Regenerar arquivos pb

```bash
protoc --go_out=. --go-grpc_out=. \
  internal/infra/grpc/protofiles/order.proto
```

#### 3B.3 Implementar ListOrders no service

**Arquivo:** `internal/infra/grpc/service/order_service.go`

```go
func (s *OrderGrpcService) ListOrders(
    ctx context.Context,
    in *pb.ListOrdersRequest,
) (*pb.ListOrdersResponse, error) {
    output, err := s.ListOrdersUseCase.Execute(usecase.ListOrdersInputDTO{})
    if err != nil {
        return nil, err
    }
    var orders []*pb.CreateOrderResponse
    for _, o := range output.Orders {
        orders = append(orders, &pb.CreateOrderResponse{
            Id:         o.ID,
            Price:      float32(o.Price),
            Tax:        float32(o.Tax),
            FinalPrice: float32(o.FinalPrice),
        })
    }
    return &pb.ListOrdersResponse{Orders: orders}, nil
}
```

**Criterio de aceite da Fase 3B:**
- [ ] Proto compila sem erros
- [ ] `go build ./internal/infra/grpc/...` passa
- [ ] gRPC server responde ao ListOrders (validar via `grpcurl`)

```bash
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
```

---

### FASE 3C — Interface GraphQL: Query listOrders

**Agente:** `golang-pro` | **Modo:** foreground | **Paralelo com 3A e 3B**

#### 3C.1 Atualizar schema

**Arquivo:** `internal/infra/graph/schema.graphqls`

```graphql
type Order {
  id: String!
  price: Float!
  tax: Float!
  finalPrice: Float!
}

input OrderInput {
  id: String!
  price: Float!
  tax: Float!
}

type Mutation {
  createOrder(input: OrderInput!): Order!
}

type Query {
  listOrders: [Order!]!  # NOVO
}
```

#### 3C.2 Regenerar codigo gqlgen

```bash
go run github.com/99designs/gqlgen generate
```

#### 3C.3 Implementar resolver

**Arquivo:** `internal/infra/graph/schema.resolvers.go`

```go
// Resolver struct precisa ter ListOrdersUseCase injetado
func (r *queryResolver) ListOrders(ctx context.Context) ([]*model.Order, error) {
    output, err := r.ListOrdersUseCase.Execute(usecase.ListOrdersInputDTO{})
    if err != nil {
        return nil, err
    }
    var orders []*model.Order
    for _, o := range output.Orders {
        orders = append(orders, &model.Order{
            ID:         o.ID,
            Price:      o.Price,
            Tax:        o.Tax,
            FinalPrice: o.FinalPrice,
        })
    }
    return orders, nil
}
```

#### 3C.4 Adicionar ListOrdersUseCase ao Resolver

**Arquivo:** `internal/infra/graph/resolver.go`

```go
type Resolver struct {
    CreateOrderUseCase usecase.CreateOrderUseCase
    ListOrdersUseCase  usecase.ListOrdersUseCase  // NOVO
}
```

**Criterio de aceite da Fase 3C:**
- [ ] `go run github.com/99designs/gqlgen generate` sem erros
- [ ] `go build ./internal/infra/graph/...` passa
- [ ] GraphQL playground em `http://localhost:8080` executa `{ listOrders { id price tax finalPrice } }` e retorna array

---

### FASE 4 — Wire DI Update

**Agente:** `golang-pro` | **Modo:** foreground | **Apos conclusao de 3A + 3B + 3C**

#### 4.1 Atualizar wire.go

**Arquivo:** `cmd/ordersystem/wire.go`

```go
//go:build wireinject
// +build wireinject

package main

import (
    "database/sql"
    "github.com/google/wire"
    // ... imports
)

var setOrderRepositoryDependency = wire.NewSet(
    database.NewOrderRepository,
    wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setCreateOrderUseCase = wire.NewSet(
    usecase.NewCreateOrderUseCase,
)

// NOVO
var setListOrdersUseCase = wire.NewSet(
    usecase.NewListOrdersUseCase,
)

func InitializeCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) CreateOrderUseCase {
    wire.Build(setOrderRepositoryDependency, setCreateOrderUseCase)
    return CreateOrderUseCase{}
}

// NOVO
func InitializeListOrdersUseCase(db *sql.DB) ListOrdersUseCase {
    wire.Build(setOrderRepositoryDependency, setListOrdersUseCase)
    return ListOrdersUseCase{}
}
```

#### 4.2 Regenerar wire_gen.go

```bash
cd cmd/ordersystem && wire
```

#### 4.3 Atualizar main.go

Injetar `ListOrdersUseCase` nos tres servidores:

```go
// REST
listOrderHandler := NewWebListOrderHandler(listOrdersUseCase)
webserver.AddHandler("GET", "/order", listOrderHandler.Create)

// gRPC — passar listOrdersUseCase ao service
grpcServer := grpc.NewOrderGrpcService(createOrderUseCase, listOrdersUseCase)

// GraphQL
resolver := &graph.Resolver{
    CreateOrderUseCase: createOrderUseCase,
    ListOrdersUseCase:  listOrdersUseCase,
}
```

**Criterio de aceite da Fase 4:**
- [ ] `wire` gera `wire_gen.go` sem erros
- [ ] `go build ./cmd/ordersystem/` compila sem erros

---

### FASE 5 — Docker Automation

**Agente:** `golang-pro` | **Modo:** foreground

#### 5.1 Dockerfile multi-stage

**Arquivo:** `Dockerfile`

```dockerfile
# Stage 1: Builder
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/ordersystem/

# Stage 2: Runner
FROM alpine:3.19

RUN apk add --no-cache mysql-client

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts/entrypoint.sh .

RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
```

#### 5.2 Script de entrypoint

**Arquivo:** `scripts/entrypoint.sh`

```bash
#!/bin/sh
set -e

echo "Aguardando MySQL ficar disponivel..."
until mysqladmin ping -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASSWORD" --silent 2>/dev/null; do
    echo "MySQL nao disponivel - aguardando..."
    sleep 2
done
echo "MySQL disponivel!"

echo "Executando migrations..."
mysql -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASSWORD" < /app/migrations/001_create_orders.sql
echo "Migrations concluidas!"

echo "Iniciando aplicacao..."
exec /app/server
```

#### 5.3 Migration SQL

**Arquivo:** `migrations/001_create_orders.sql`

```sql
CREATE DATABASE IF NOT EXISTS orders;
USE orders;

CREATE TABLE IF NOT EXISTS orders (
    id          VARCHAR(255) NOT NULL PRIMARY KEY,
    price       DOUBLE       NOT NULL,
    tax         DOUBLE       NOT NULL,
    final_price DOUBLE       NOT NULL
);
```

#### 5.4 docker-compose.yaml completo

**Arquivo:** `docker-compose.yaml`

```yaml
version: "3.8"

services:
  mysql:
    image: mysql:8.0
    container_name: mysql
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: orders
    ports:
      - "3306:3306"
    volumes:
      - .docker/mysql:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-proot"]
      interval: 5s
      timeout: 5s
      retries: 10
      start_period: 10s

  rabbitmq:
    image: rabbitmq:3-management
    container_name: rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    healthcheck:
      test: ["CMD", "rabbitmqctl", "status"]
      interval: 10s
      timeout: 5s
      retries: 5

  app:
    build: .
    container_name: clean-arch-app
    ports:
      - "8000:8000"    # REST
      - "50051:50051"  # gRPC
      - "8080:8080"    # GraphQL
    environment:
      DB_DRIVER: mysql
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: root
      DB_PASSWORD: root
      DB_NAME: orders
      WEB_SERVER_PORT: :8000
      GRPC_SERVER_PORT: 50051
      GRAPHQL_SERVER_PORT: 8080
      RABBITMQ_DSN: amqp://guest:guest@rabbitmq:5672/
    depends_on:
      mysql:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
    restart: on-failure
```

**Criterio de aceite da Fase 5:**
- [ ] `docker compose up --build` completa sem erros
- [ ] `docker compose ps` mostra todos os servicos como `healthy` ou `running`
- [ ] Logs do `app` mostram "MySQL disponivel!", "Migrations concluidas!", "Iniciando aplicacao..."
- [ ] Tabela `orders` existe no MySQL apos o startup

---

### FASE 6 — api.http e README

**Agente:** `golang-pro` | **Modo:** foreground

#### 6.1 api.http completo

**Arquivo:** `api.http`

```http
### ==========================================
### CREATE ORDER — REST
### ==========================================
POST http://localhost:8000/order
Content-Type: application/json

{
  "id": "order-001",
  "price": 100.50,
  "tax": 10.50
}

###

POST http://localhost:8000/order
Content-Type: application/json

{
  "id": "order-002",
  "price": 200.00,
  "tax": 20.00
}

### ==========================================
### LIST ORDERS — REST
### ==========================================
GET http://localhost:8000/order

### ==========================================
### CREATE ORDER — GraphQL
### ==========================================
POST http://localhost:8080/query
Content-Type: application/json

{
  "query": "mutation { createOrder(input: { id: \"order-003\", price: 50.0, tax: 5.0 }) { id price tax finalPrice } }"
}

### ==========================================
### LIST ORDERS — GraphQL
### ==========================================
POST http://localhost:8080/query
Content-Type: application/json

{
  "query": "{ listOrders { id price tax finalPrice } }"
}
```

#### 6.2 README.md

**Arquivo:** `README.md`

```markdown
# Clean Architecture — Go

Use case ListOrders exposto simultaneamente via REST, gRPC e GraphQL.

## Execucao

docker compose up --build

## Portas

| Protocolo | Porta | Endpoint |
|-----------|-------|----------|
| REST      | 8000  | GET /order, POST /order |
| gRPC      | 50051 | pb.OrderService/ListOrders |
| GraphQL   | 8080  | POST /query |

## Testando

### REST
curl http://localhost:8000/order

### GraphQL
Acesse http://localhost:8080 para o playground interativo.

### gRPC
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
```

---

### FASE 7 — Revisao, Seguranca e Testes (Paralelo)

#### 7A — code-reviewer

**Agente:** `code-reviewer` | **Modo:** background

Focos de revisao:
- Fronteiras de camadas respeitadas (domain nao importa infra)
- Injecao de dependencia via interface (nao implementacao concreta)
- Nomes de pacotes e arquivos em snake_case
- Ausencia de logica de negocio nos handlers (REST/gRPC/GraphQL)
- Wire DI corretamente configurado

#### 7B — security-auditor

**Agente:** `security-auditor` | **Modo:** background

Focos:
- SQL injection em `GetAll()` e `Save()`
- Credenciais hardcoded no codigo (nao apenas no .env)
- Portas expostas desnecessariamente
- Input validation na camada web

#### 7C — test-runner

**Agente:** `test-runner` | **Modo:** background

Comandos de validacao:

```bash
# Testes unitarios
go test ./... -v

# Build completo
go build ./...

# Docker end-to-end
docker compose up --build -d
sleep 15
curl -s http://localhost:8000/order
docker compose down
```

---

## 7. Criterios de Aceite Globais

- [ ] `go build ./cmd/ordersystem/` compila sem erros ou warnings
- [ ] `go test ./...` retorna apenas PASS (zero FAIL)
- [ ] `docker compose up --build` conclui sem erros
- [ ] `docker compose ps` mostra todos os servicos saudaveis
- [ ] `curl -s http://localhost:8000/order` retorna JSON (array)
- [ ] `curl -s -X POST http://localhost:8000/order -d '{"id":"x","price":10,"tax":1}'` retorna order criada
- [ ] GraphQL query `{ listOrders { id price tax finalPrice } }` retorna array
- [ ] gRPC `grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders` retorna resposta
- [ ] Migrations executam automaticamente (tabela `orders` existe apos startup)
- [ ] MySQL healthcheck passa antes do app iniciar (sem race condition)
- [ ] Nenhuma credencial hardcoded no codigo-fonte
- [ ] Fronteiras de Clean Architecture respeitadas (domain nao depende de infra)

---

## 8. Riscos e Mitigacoes

| Risco | Probabilidade | Mitigacao |
|-------|--------------|-----------|
| Race condition MySQL startup | ALTA | entrypoint.sh com loop de retry |
| protoc nao instalado | MEDIA | Incluir no Dockerfile stage builder |
| wire regeneracao manual errada | MEDIA | `wire` rodado automaticamente no Makefile |
| gqlgen breaking changes | BAIXA | Versao fixada no go.mod |
| graphql resolver com campo errado (finalPrice vs final_price) | MEDIA | Validar schema vs model gerado |

---

## 9. Perguntas em Aberto (Bloqueios Criticos)

Antes de iniciar a implementacao, confirmar:

1. **Modulo Go**: usar `github.com/reinaldosaraiva/clean-arch`?
2. **Banco de dados**: MySQL (como no template de referencia) ou PostgreSQL?
3. **Autenticacao**: fora do escopo desta iteracao?
4. **OrderItem**: entidade `Order` tem apenas os 4 campos do template (id, price, tax, final_price)?
5. **Repositorio GitHub**: ja criado e vazio, ou criar durante a implementacao?

---

## 10. Prompt para o Multi-Agent-System

```xml
<objetivo>
Implementar Clean Architecture em Go com o use case ListOrders exposto
simultaneamente via REST (GET /order), gRPC (OrderService.ListOrders) e
GraphQL (query listOrders). Docker completamente automatizado:
docker compose up sobe MySQL, executa migrations e inicia a aplicacao.
</objetivo>

<restricoes>
- NAO inventar arquivos sem cria-los explicitamente
- Citar arquivos reais do repositorio antes de qualquer escrita de codigo
- Exigir evidencia verificavel: diff por arquivo + output de go test + docker compose up
- NAO declarar "pronto" sem checklist de aceite totalmente verificado
- NAO usar libs sem adiciona-las ao go.mod via go get
- Cada agente entrega apenas sua fatia; nao declara conclusao sem evidencia de build
- Fronteiras Clean Architecture DEVEM ser respeitadas (domain nao importa infra)
</restricoes>

<escopo>
1. Modulo Go: github.com/reinaldosaraiva/clean-arch
2. Entidade Order: ID, Price, Tax, FinalPrice
3. UseCase ListOrders: GetAll() no repositorio MySQL
4. REST: GET /order retorna JSON array
5. gRPC: OrderService.ListOrders retorna repeated Order
6. GraphQL: query listOrders retorna [Order!]!
7. Dockerfile multi-stage + docker-compose com healthcheck
8. entrypoint.sh: wait-for-mysql + migrations + start
9. Testes unitarios no usecase e entidade
</escopo>

<formato-de-resposta>
Cada agente DEVE responder com:
1. ARQUIVOS ENCONTRADOS (Glob/Read antes de qualquer acao)
2. PLANO DE IMPLEMENTACAO (passos sequenciais com dependencias)
3. ALTERACOES (path absoluto + diff ou conteudo completo + justificativa)
4. TESTES EXECUTADOS (comando + output completo)
5. LACUNAS E PROXIMOS PASSOS (o que nao foi feito + bloqueantes)
</formato-de-resposta>

<criterios-de-aceite>
- go build ./cmd/ordersystem/ compila sem erros
- go test ./... retorna apenas PASS
- curl http://localhost:8000/order retorna JSON
- GraphQL query listOrders retorna array
- grpcurl ListOrders retorna resposta
- docker compose up --build sem erros
- Tabela orders existe apos startup automatico
</criterios-de-aceite>
```

---

*Plano gerado em: 2026-04-02*
*Score de complexidade: 7 (COMPLEX)*
*Modelo recomendado: opus*
