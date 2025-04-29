FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .

RUN go mod init order-service && \
    go get -u github.com/gorilla/mux && \
    go get -u github.com/joho/godotenv && \
    go build -o order-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/order-service .
COPY .env .


EXPOSE 8082
CMD ["./order-service"]
