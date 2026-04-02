# Stage 1: Builder
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/ordersystem/

# Stage 2: Runner
FROM alpine:3.19

RUN apk add --no-cache mysql-client

WORKDIR /app

COPY --from=builder /app/server .
COPY migrations/ ./migrations/
COPY scripts/entrypoint.sh ./entrypoint.sh

RUN chmod +x entrypoint.sh

EXPOSE 8000 50051 8080

ENTRYPOINT ["./entrypoint.sh"]
