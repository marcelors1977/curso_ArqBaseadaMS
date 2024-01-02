FROM golang:1.18-alpine AS builder

ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY . .

RUN apk update && apk add gcc librdkafka-dev zstd-libs libsasl lz4-dev libc-dev musl-dev 

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -tags musl -o main ./cmd/walletcore/main.go

FROM alpine

RUN apk --no-cache upgrade && apk add librdkafka

COPY --from=builder /app/main /


CMD ["tail", "-f", "/dev/null"]