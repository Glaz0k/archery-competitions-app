FROM golang:1.23.3 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app-server ./cmd/service/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app-server .

EXPOSE 8080
ENTRYPOINT ["./app-server"]