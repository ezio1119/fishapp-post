# 開発用
FROM golang:1.13-alpine AS dev

WORKDIR /app
RUN apk add --no-cache tzdata git && \
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

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY --from=builder /src/main .
COPY --from=builder /src/conf/conf.yml /app/conf/conf.yml

CMD ["./main"]