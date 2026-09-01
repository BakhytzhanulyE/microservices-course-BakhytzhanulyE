# ============================
# Этап 1: сборка
# ============================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# GOWORK=off: внутри образа лежат только три нужных модуля, а go.work перечисляет все
# восемь — в workspace-режиме сборка упала бы на отсутствующих каталогах.
# Связь между модулями обеспечивают replace-директивы в go.mod.
ENV GOWORK=off

# Общие модули копируем первыми: они меняются реже кода сервиса,
# поэтому слой с зависимостями переиспользуется между сборками.
COPY platform/ ./platform/
COPY shared/ ./shared/
COPY inventory/go.mod inventory/go.sum ./inventory/

WORKDIR /app/inventory
RUN go mod download

COPY inventory/ ./

# grpc-health-probe умеет дергать стандартный grpc.health.v1 — им healthcheck
# проверяет живость сервиса, у которого нет HTTP-порта.
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.38/grpc_health_probe-linux-amd64 /app/bin/grpc-health-probe
RUN chmod +x /app/bin/grpc-health-probe

# CGO_ENABLED=0 даёт статический бинарник — он работает в пустом alpine без libc.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/bin/inventory ./cmd/inventory

# ============================
# Этап 2: рабочий образ
# ============================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates wget && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/inventory ./inventory
COPY --from=builder /app/bin/grpc-health-probe /bin/grpc-health-probe

USER appuser

ENTRYPOINT ["./inventory"]
