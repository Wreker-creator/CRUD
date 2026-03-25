FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o taskmanager ./cmd/webserver

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/taskmanager .

EXPOSE 5001

CMD ["./taskmanager"]