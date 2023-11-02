FROM golang:1.20 AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./...
FROM alpine:3.14 AS prodction
COPY --from=builder /app .
CMD ["./main"]
