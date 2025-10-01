FROM golang:1.22 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o orders-service cmd/main.go

FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=builder /app/orders-service .
COPY --from=builder /app/internal/postgres/model.sql ./model.sql
COPY --from=builder /app/config.env .
COPY --from=builder /app/web/index.html .
CMD ["./orders-service"]