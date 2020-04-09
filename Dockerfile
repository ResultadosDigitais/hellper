FROM golang:1.13-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hellper /app/cmd/http

FROM alpine:latest
COPY --from=builder /app/hellper /app/hellper
EXPOSE 8080

ENTRYPOINT ["/app/hellper"]
