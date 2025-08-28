FROM golang:1.22 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o orders-service

FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=builder /app/orders-service .
COPY --from=builder /app/postgres/model.sql ./model.sql
CMD ["./orders-service"]