FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .

RUN go mod init restaurant-service && \
    go get -u github.com/gorilla/mux && \
    go build -o restaurant-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/restaurant-service .

EXPOSE 8081
CMD ["./restaurant-service"]