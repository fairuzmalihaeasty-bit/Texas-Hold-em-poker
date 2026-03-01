### Multi-stage build: produce a static binary and package it in a tiny image
FROM golang:1.20-alpine AS builder
WORKDIR /src

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy project
COPY . .

# Build server binary
RUN go build -trimpath -o /texasholdem ./cmd/server

# Final image
FROM scratch
COPY --from=builder /texasholdem /texasholdem
COPY --from=builder /src/web /web

EXPOSE 8080
ENTRYPOINT ["/texasholdem"]
