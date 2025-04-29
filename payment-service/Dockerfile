FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .

RUN go mod init payment-service && \
    go get -u github.com/gorilla/mux && \
    go get -u github.com/joho/godotenv && \
    go build -o payment-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/payment-service .
COPY .env .


EXPOSE 8083
CMD ["./payment-service"]