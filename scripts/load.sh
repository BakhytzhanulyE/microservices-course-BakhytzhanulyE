#!/bin/sh
# Нагрузочный прогон для поиска гонок в order.
# Смысл: параллельно бьём в запись (POST) и в чтение (GET) одного хранилища,
# чтобы -race успел заметить незащищённый доступ.
#
# Запуск (сервер должен быть поднят отдельно, обязательно с -race):
#   go run -race ./order/cmd/order      # в соседней вкладке
#   sh scripts/load.sh                  # здесь
#
# Параметры через переменные окружения:
#   N=400 P=100 HOST=localhost:8080 sh scripts/load.sh

set -eu

N="${N:-400}"          # сколько итераций
P="${P:-100}"          # сколько параллельных процессов
HOST="${HOST:-localhost:8080}"

if ! curl -s -o /dev/null --max-time 2 "http://$HOST/health"; then
  echo "❌ $HOST не отвечает на /health — подними сервер: go run -race ./order/cmd/order" >&2
  exit 1
fi

echo "🔥 $N итераций, параллельность $P, цель $HOST"

seq 1 "$N" | xargs -P "$P" -I{} sh -c '
  curl -s -o /dev/null -X POST "http://$0/api/v1/orders" \
    -d "{\"user_uuid\":\"u$1\",\"total_price\":1}" &
  curl -s -o /dev/null "http://$0/api/v1/orders/1"
  wait
' "$HOST" {}

echo "✅ Прогон закончен. Смотри вкладку с сервером: есть ли там WARNING: DATA RACE"
