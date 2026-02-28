# ---------- Build stage ----------
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o user-service ./cmd/user/main.go

# ---------- Runtime stage ----------
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/user-service .

EXPOSE 7071

CMD ["./user-service"]


