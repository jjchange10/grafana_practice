FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

CMD ["/app/main"]
