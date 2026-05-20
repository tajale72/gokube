# Build stage
FROM golang:1.24 AS builder  
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/gokube

# Final stage
FROM alpine:3.20
RUN apk update && apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./main"]
