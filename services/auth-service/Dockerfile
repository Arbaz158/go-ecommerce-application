# ---------- Build stage ----------
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o auth-service ./cmd/auth/main.go

# ---------- Runtime stage ----------
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/auth-service .
EXPOSE 7070

CMD ["./auth-service"]