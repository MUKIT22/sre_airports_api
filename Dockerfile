# Build stage
FROM golang:1.23.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bd-airports

# Final stage
FROM gcr.io/distroless/base-debian10
WORKDIR /root/
COPY --from=builder /app/bd-airports .
EXPOSE 8080
CMD ["./bd-airports"]

