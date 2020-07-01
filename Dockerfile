# 開発用
FROM golang:1.13-alpine AS dev

WORKDIR /app
RUN apk add --no-cache tzdata git && \
    go get github.com/pilu/fresh
# skaffold sync用
COPY . .

CMD ["fresh"]

# コンパイラ用
FROM golang:1.13-alpine AS builder1
WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main .

# grpc_health_probe
FROM busybox as builder2
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# 本番用
FROM alpine as prod
WORKDIR /app

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY --from=builder1 /src/main .
COPY --from=builder2 /bin/grpc_health_probe /bin/grpc_health_probe
COPY conf/conf.yml /app/conf/conf.yml

CMD ["./main"]