# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependencies separately from source
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o somq ./cmd/main.go

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.20 AS runtime

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/somq .
COPY config/config.yml config/config.yml

ENV CONFIG_PATH=config/config.yml

EXPOSE 8090

ENTRYPOINT ["./somq"]
