FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o orders-service cmd/main.go

FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=builder /app/orders-service .
COPY --from=builder /app/internal/infrastructure/postgres/migrations/model_up.sql ./migrations/model_up.sql
COPY --from=builder /app/internal/infrastructure/postgres/migrations/model_down.sql ./migrations/model_down.sql
COPY --from=builder /app/config.env .
COPY --from=builder /app/web/index.html .
CMD ["./orders-service"]