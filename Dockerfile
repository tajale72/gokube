# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod .
COPY cmd/ cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/gokube

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"] 