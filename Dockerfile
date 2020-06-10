# 開発用
FROM golang:1.13-alpine AS dev

WORKDIR /app
RUN apk add --no-cache tzdata git gcc && \
    go get github.com/pilu/fresh

CMD ["fresh"]

# コンパイラ用
FROM golang:1.13-alpine AS builder
WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main .

# 本番用
FROM alpine as prod
WORKDIR /app
RUN apk add --no-cache tzdata

COPY --from=builder /src/main .
COPY --from=builder /src/conf/conf.yml /app/conf/conf.yml

CMD ["./main"]