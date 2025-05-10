FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/jwt-auth-service ./cmd/jwt-auth-service/main.go

FROM alpine:3.19
RUN apk add --no-cache bash

COPY --from=builder /app/config /config
COPY --from=builder /app/.env .env
COPY --from=builder /build/jwt-auth-service /jwt-auth-service

ENTRYPOINT ["/jwt-auth-service"]

CMD ["--config_path="]