# Development stage: used by Docker Compose for live reload.
FROM golang:1.26-alpine AS development

WORKDIR /app

RUN go install github.com/air-verse/air@v1.67.3
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]

# Test stage: used to run unit tests.
FROM golang:1.26-alpine AS test

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["go", "test", "./internal/...", "-v"]

# Build stage: used for the production image.
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /server \
    ./cmd/server

# Runtime stage
FROM alpine:3.23 AS production

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /server /app/server

EXPOSE 8080

CMD ["/app/server"]
