FROM golang:1.26-alpine AS builder
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.23
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/goose /usr/local/bin/goose
ENTRYPOINT ["goose"]