# Build stage
FROM golang:alpine AS builder
# SQLite requires gcc and musl-dev to compile CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go app with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o offline-sync main.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/offline-sync .
# Create data directory for SQLite database
RUN mkdir -p /app/data 

EXPOSE 8000

CMD ["./offline-sync"]