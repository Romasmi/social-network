FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations/ ./migrations/
COPY --from=builder /app/*.yaml ./

EXPOSE 8000

CMD ["./main"]