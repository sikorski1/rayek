# Etap 1 — build
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server .

# Etap 2 — minimalny obraz do uruchomienia
FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/server .

COPY --from=builder /app/data ./data

ENV GIN_MODE=release
EXPOSE 3000
CMD ["./server"]