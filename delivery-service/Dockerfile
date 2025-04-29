FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .

RUN go mod init delivery-service && \
    go get -u github.com/gorilla/mux && \
    go get -u github.com/joho/godotenv && \
    go build -o delivery-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/delivery-service .
COPY .env . 


EXPOSE 8084
CMD ["./delivery-service"]