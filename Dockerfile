FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    pkg-config \
    gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# IMPORTANT: disable vendored librdkafka
ENV CGO_ENABLED=1
ENV GOFLAGS="-tags=dynamic"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api cmd/api/main.go
RUN go build -o worker cmd/worker/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/worker .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/*.yaml ./

EXPOSE 8888
CMD ["./api"]
