FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .

RUN go mod init notification-service && \
    go get -u github.com/gorilla/mux && \
    go build -o notification-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/notification-service .

EXPOSE 8085
CMD ["./notification-service"]