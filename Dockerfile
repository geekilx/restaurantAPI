FROM golang:1.25.4-alpine AS builder

WORKDIR /build

RUN go install github.com/air-verse/air@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o restaurant ./cmd/api

# second stage
FROM scratch


COPY --from=builder /build/restaurant .

# Start the application
CMD ["./restaurant"]
