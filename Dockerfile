FROM golang:1.20 AS builder

WORKDIR /app
COPY . .

RUN go build -o bin/loadbalancer main.go

FROM alpine:3.16

WORKDIR /root/
COPY --from=builder /app/cmd/loadbalancer .

EXPOSE 8080

CMD ["./loadbalancer"]
