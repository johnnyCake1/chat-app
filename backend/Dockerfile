FROM golang:1.22.1-alpine

RUN apk add --no-cache git gcc musl-dev bash

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/app/

# Copy the wait-for-it.sh script and make it executable
COPY scripts/wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

# This is to ensure that the database and RabbitMQ will always go up first before running the server
CMD ["./wait-for-it.sh", "postgres:5432", "--", "./wait-for-it.sh", "rabbitmq:5672", "--timeout=30", "--", "./main"]
