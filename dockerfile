FROM golang:1.24 AS builder
WORKDIR /app
COPY . .

RUN go build -o api ./cmd/api && \
    go build -o migrate ./cmd/migrate

FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /app/api /app/api
COPY --from=builder /app/migrate /app/migrate

