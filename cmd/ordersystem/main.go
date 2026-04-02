package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/reinaldosaraiva/clean-arch/configs"
	"github.com/reinaldosaraiva/clean-arch/internal/event"
	"github.com/reinaldosaraiva/clean-arch/internal/event/handler"
	"github.com/reinaldosaraiva/clean-arch/internal/infra/database"
	"github.com/reinaldosaraiva/clean-arch/internal/infra/graph"
	grpcService "github.com/reinaldosaraiva/clean-arch/internal/infra/grpc/service"
	"github.com/reinaldosaraiva/clean-arch/internal/infra/web"
	"github.com/reinaldosaraiva/clean-arch/internal/infra/web/webserver"
	"github.com/reinaldosaraiva/clean-arch/internal/usecase"
	"github.com/reinaldosaraiva/clean-arch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/reinaldosaraiva/clean-arch/internal/infra/grpc/pb"
)

func main() {
	// Config
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// MySQL
	db, err := sql.Open(cfg.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// RabbitMQ
	rabbitConn, err := amqp.Dial(cfg.RabbitMQDSN)
	if err != nil {
		panic(err)
	}
	defer rabbitConn.Close()

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		panic(err)
	}
	defer rabbitCh.Close()

	// Events
	orderCreatedEvent := event.NewOrderCreated()
	eventDispatcher := events.NewEventDispatcher()
	orderCreatedHandler := handler.NewOrderCreatedHandler(rabbitCh)
	eventDispatcher.Register("OrderCreated", orderCreatedHandler)

	// Repository
	orderRepo := database.NewOrderRepository(db)

	// Use cases
	createOrderUseCase := *usecase.NewCreateOrderUseCase(orderRepo, orderCreatedEvent, eventDispatcher)
	listOrdersUseCase := *usecase.NewListOrdersUseCase(orderRepo)

	// --- REST Server ---
	createOrderHandler := web.NewWebOrderHandler(createOrderUseCase)
	listOrderHandler := web.NewWebListOrderHandler(listOrdersUseCase)

	ws := webserver.NewWebServer(cfg.WebServerPort)
	ws.AddHandler("POST", "/order", createOrderHandler.Create)
	ws.AddHandler("GET", "/order", listOrderHandler.List)

	// --- gRPC Server ---
	grpcSrv := grpc.NewServer()
	orderGrpcService := grpcService.NewOrderGrpcService(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcSrv, orderGrpcService)
	reflection.Register(grpcSrv)

	// --- GraphQL Server ---
	resolver := &graph.Resolver{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
	gqlSrv := graphqlHandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Start all 3 servers
	fmt.Println("Starting REST server on", cfg.WebServerPort)
	fmt.Printf("Starting gRPC server on :%d\n", cfg.GRPCServerPort)
	fmt.Printf("Starting GraphQL server on :%d\n", cfg.GraphQLServerPort)

	// gRPC in goroutine
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCServerPort))
		if err != nil {
			panic(err)
		}
		if err := grpcSrv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// GraphQL in goroutine
	go func() {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Handle("/", playground.Handler("GraphQL Playground", "/query"))
		r.Handle("/query", gqlSrv)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GraphQLServerPort), r); err != nil {
			panic(err)
		}
	}()

	// REST on main goroutine (blocking)
	ws.Start()
}
