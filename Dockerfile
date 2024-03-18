FROM golang:1.22.1-alpine3.18 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o darth-veda .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/darth-veda .

EXPOSE 8080

CMD ["./darth-veda", "--mode=server"]
