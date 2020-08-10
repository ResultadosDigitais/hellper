FROM golang:1.14-buster AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o . ./cmd/http ./cmd/notify

FROM debian:buster
RUN apt update -y && apt upgrade -y && apt install ca-certificates -y
COPY --from=builder /app/entrypoint.sh /app/http /app/notify /app/
EXPOSE 8080

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]
