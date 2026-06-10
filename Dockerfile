# ── base ────────────────────────────────────────────────────────
FROM golang:alpine AS base

ARG BUILD_ENV=production
ENV BUILD_ENV=${BUILD_ENV}

WORKDIR /app

# ── installer ───────────────────────────────────────────────────
FROM base AS installer

COPY go.mod go.sum ./
RUN go mod download

# ── builder ─────────────────────────────────────────────────────
FROM base AS builder

COPY --from=installer /go/pkg /go/pkg
COPY --from=installer /app/go.mod /app/go.sum ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o server .

# ── release ─────────────────────────────────────────────────────
FROM alpine:latest AS release

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3000
ENTRYPOINT ["./server"]
