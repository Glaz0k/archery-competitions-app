FROM golang:1.23.3 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app-server ./cmd/service/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app-migrate ./cmd/migrate/main.go

FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app
COPY --from=builder /app-server .
COPY --from=builder /app-migrate .
COPY cmd/migrate/migrations ./cmd/migrate/migrations
COPY entrypoint.sh .

RUN chmod +x entrypoint.sh && dos2unix entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/bin/bash", "./entrypoint.sh"]