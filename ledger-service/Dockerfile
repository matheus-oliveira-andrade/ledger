# Build stage
FROM golang:alpine3.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Final stage
FROM alpine:3.19

# Install ca-certificates for any external API calls
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

RUN addgroup -g 1001 appuser && adduser -D -u 1001 -G appuser appuser

COPY --from=builder /build/main .

RUN chown appuser:appuser main

USER appuser

EXPOSE 8080

ENTRYPOINT ["./main"]