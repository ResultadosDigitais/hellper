FROM golang:1.14-buster AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hellper /app/cmd/http

# POC Notify - build notify
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notify /app/cmd/notify
RUN sed -i 's/^\(.*\)$/export \1/g' /app/.env

FROM debian:buster
RUN apt update -y && apt upgrade -y && apt install procps cron ca-certificates -y
COPY --from=builder /app/hellper /app/hellper

# POC Notify - cmd, cron, env
COPY --from=builder /app/.env /app/
COPY --from=builder /app/notify /app/notify
COPY --from=builder /app/scripts/cron/*-cron /etc/cron.d/
RUN chmod 0600 /etc/cron.d/*-cron

EXPOSE 8080

ENTRYPOINT cron start && /app/hellper
