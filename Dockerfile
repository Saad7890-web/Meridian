# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o meridian ./cmd/meridian

# Run stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/meridian .

EXPOSE 8080

CMD ["./meridian"]