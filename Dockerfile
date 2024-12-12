FROM golang:1.23 as builder

WORKDIR /usr/src/app
RUN go mod download
COPY . .
RUN go build -v -o wallets ./...

FROM debian:stable-slim
WORKDIR /bin
COPY --from=builder /usr/src/app/wallets ./
ENTRYPOINT ["./wallets"]
