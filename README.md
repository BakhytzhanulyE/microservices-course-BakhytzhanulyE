# microservices-course-BakhytzhanulyE

[![CI](https://github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/actions/workflows/ci.yml/badge.svg)](https://github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/actions/workflows/ci.yml)

Система из семи Go-модулей: заказ деталей для космического корабля, оплата, сборка и уведомления.
Сервисы общаются двумя способами — синхронно по gRPC там, где нужен ответ прямо сейчас,
и асинхронно через Kafka там, где ответ не нужен.

## Как это работает

```
                    ┌──────────────┐
   HTTP :8080  ───▶ │    order     │ ── gRPC ─▶ inventory :50051 ── MongoDB
                    │  PostgreSQL  │ ── gRPC ─▶ payment   :50052
                    └──────┬───────┘
                           │ publish: order.paid
                           ▼
                    ┌──────────────┐
                    │    Kafka     │
                    └──┬────────┬──┘
        order.paid     │        │     order.paid
                       ▼        ▼     ship.assembled
                ┌──────────┐  ┌──────────────┐
                │ assembly │  │ notification │ ──▶ Telegram
                └────┬─────┘  └──────────────┘
                     │ publish: ship.assembled
                     └──────────▶ Kafka

   iam :50053 ── PostgreSQL (пользователи) + Redis (сессии)
```

Путь одного заказа:

1. `POST /api/v1/orders` — `order` спрашивает у `inventory` цены деталей, считает сумму,
   сохраняет заказ в PostgreSQL со статусом `PENDING_PAYMENT`.
2. `POST /api/v1/orders/{uuid}/pay` — `order` вызывает `payment` по gRPC, сохраняет
   UUID транзакции и публикует событие `order.paid`.
3. `assembly` читает `order.paid`, «собирает корабль» и публикует `ship.assembled`.
4. `notification` читает оба топика и шлёт сообщения в Telegram.

## Модули

Каждый сервис — отдельный Go-модуль со своим `go.mod`. Связь между ними — через
`replace`-директивы и `go.work`.

| Модуль         | Что делает                                     | Транспорт   | Хранилище            |
|----------------|------------------------------------------------|-------------|----------------------|
| `order`        | оформление и оплата заказов, оркестрация        | HTTP :8080  | PostgreSQL           |
| `inventory`    | каталог деталей                                 | gRPC :50051 | MongoDB              |
| `payment`      | проведение оплаты                               | gRPC :50052 | —                    |
| `iam`          | регистрация, вход, JWT                          | gRPC :50053 | PostgreSQL + Redis   |
| `assembly`     | сборка корабля по оплаченному заказу            | Kafka       | —                    |
| `notification` | уведомления в Telegram                          | Kafka       | —                    |
| `platform`     | общий код: логгер, closer, Kafka, кеш, метрики  | —           | —                    |
| `shared`       | контракты: `.proto`, сгенерированный код, OpenAPI | —         | —                    |

### Слои внутри сервиса

Одинаковые во всех сервисах, зависимости смотрят только внутрь:

```
cmd/          точка входа: загрузка конфига и сигналы
internal/
  app/        сборка приложения (app.go) и контейнер зависимостей (di.go)
  config/     разбор переменных окружения в типизированные интерфейсы
  api/        транспорт: gRPC-хендлеры или HTTP-ручки
  service/    бизнес-логика, ничего не знает про транспорт и базу
  repository/ работа с хранилищем
  converter/  перевод между моделями слоёв
  model/      доменные сущности
```

## Стек

- **Go 1.26**
- **gRPC** + **protobuf**, контракты собираются через **buf**
- **HTTP** на **chi**, контракт описан в `shared/api/order/v1/order.openapi.yaml`
- **Kafka** (IBM/sarama) — событийная шина
- **PostgreSQL** (pgx/v5) + миграции **goose**, **MongoDB**, **Redis**
- **JWT** (golang-jwt) + **bcrypt** — аутентификация
- **zap** — структурированные логи, **Prometheus** — метрики, **OpenTelemetry** + **Jaeger** — трейсы
- **Docker Compose** — локальный запуск, **Task** — команды проекта
- **golangci-lint**, **testify** — качество кода, **GitHub Actions** — CI

## Требования

- [Go](https://go.dev/dl/) 1.26 или новее
- [Task](https://taskfile.dev/installation/) — `brew install go-task`
- Docker с плагином Compose

Линтер, форматтеры, `buf`, `grpcurl` и `envsubst` ставить руками не нужно — они
скачиваются в `./bin` нужных версий автоматически.

## Быстрый старт

```bash
task env:generate     # разложить .env по сервисам из deploy/env/.env
task up-all           # поднять всю систему в Docker
task ps               # посмотреть, что запустилось
task http:order:scenario   # прогнать сценарий: заказ → оплата → сборка
task logs -- assembly      # посмотреть, как собирается корабль
task down-all         # остановить всё
```

После `task http:order:scenario` открой http://localhost:16686, выбери сервис
`order` и найди последний трейс: там должен быть один трейс на четыре сервиса —
`order`, `payment`, `assembly`, `notification`.

Первый `task env:generate` создаёт `deploy/env/.env` из шаблона. Дальше правится
только этот файл — по нему генерируются `.env` каждого сервиса. Сам `.env` в git не
попадает, шаблоны — попадают.

### Что где открывается

| Адрес                    | Что это                     |
|--------------------------|-----------------------------|
| http://localhost:8080    | HTTP API заказов            |
| http://localhost:8090    | Kafka UI                    |
| http://localhost:16686   | Jaeger — трейсы             |
| http://localhost:9091    | Prometheus                  |
| http://localhost:3000    | Grafana (admin/admin)       |

### Запуск сервиса с локальной машины

Инфраструктуру поднимаем в Docker, а сам сервис запускаем из исходников — так
удобнее отлаживать:

```bash
task up-core
task up-inventory
go run ./order/cmd/order
```

В `deploy/compose/*/.env` хосты указаны как `localhost`, а в docker-compose поверх
подставляются имена контейнеров. Поэтому один и тот же `.env` работает в обоих случаях.

## HTTP API заказов

Полный контракт — в `shared/api/order/v1/order.openapi.yaml`.

```bash
# создать заказ
curl -X POST http://localhost:8080/api/v1/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_uuid": "22222222-2222-4222-8222-222222222222", "part_uuids": ["<part_uuid>"]}'

# оплатить
curl -X POST http://localhost:8080/api/v1/orders/<order_uuid>/pay \
  -H 'Content-Type: application/json' \
  -d '{"payment_method": "CARD"}'

# посмотреть
curl http://localhost:8080/api/v1/orders/<order_uuid>

# отменить (только неоплаченный)
curl -X POST http://localhost:8080/api/v1/orders/<order_uuid>/cancel
```

Коды ответов: `400` — плохой запрос или неизвестная деталь, `404` — нет заказа,
`409` — заказ уже оплачен или отменён, `502` — упал платёжный сервис,
`503` — недоступен каталог.

## gRPC API

Все gRPC-сервисы включают reflection, поэтому `grpcurl` работает без `.proto` под рукой:

```bash
task grpc:inventory:list          # весь каталог деталей
task grpc:inventory:get -- <uuid> # одна деталь
task grpc:payment:pay             # прямая оплата
task grpc:iam:register            # регистрация пользователя
task grpc:iam:login               # вход, вернёт пару токенов
```

## Команды

```bash
task format          # gofumpt + gci по всем модулям
task lint            # golangci-lint по всем модулям
task build           # собрать все модули
task test            # прогнать тесты с -race
task test:coverage   # покрытие в coverage/coverage.out
task deps:update     # go work sync + go mod tidy во всех модулях
task proto:gen       # перегенерировать код из .proto
task --list          # весь список команд
```

Версии инструментов зафиксированы в `Taskfile.yml`, CI берёт их оттуда же
через `.github/scripts/extract-versions.sh` — локальная машина и CI не разъезжаются.

## Контракты

`.proto`-файлы лежат в `shared/proto`, сгенерированный код — в `shared/pkg/proto`.
После правки контракта:

```bash
task proto:gen
task deps:update
```

| Контракт                          | Кто отдаёт  | Кто использует             |
|-----------------------------------|-------------|----------------------------|
| `inventory/v1/inventory.proto`    | `inventory` | `order`                    |
| `payment/v1/payment.proto`        | `payment`   | `order`                    |
| `iam/v1/iam.proto`                | `iam`       | клиенты                    |
| `events/v1/order.proto`           | —           | `order`, `assembly`, `notification` |

## Тесты

Бизнес-логика покрыта unit-тестами на заглушках — без Docker и без сети:

```bash
task test
```

Что проверяется: расчёт суммы заказа и схлопывание дублей деталей, отказ при
неизвестной детали, повторная оплата (409), сохранность статуса при сбое платежа,
живучесть оплаты при упавшей Kafka, коды HTTP-ответов, регистрация и вход в `iam`
(включая протухший токен и удалённую сессию), идемпотентность `closer`,
склейка трейса через Kafka и разделение клиентских ошибок от серверных.

## Наблюдаемость

- **Логи** — zap в JSON, у HTTP-запросов есть `trace_id` из заголовка `X-Request-Id`.
  Уровень выбирается по коду ответа: отказ по бизнес-правилу (`AlreadyExists`,
  `NotFound`, `InvalidArgument` и подобные) пишется как `WARN`, а `ERROR`
  остаётся за настоящими сбоями — иначе алерты по частоте ошибок срабатывали бы
  на обычном поведении пользователей.
- **Метрики** — каждый сервис отдаёт `/metrics`, Prometheus собирает все шесть.
  В метке пути стоит шаблон маршрута (`/api/v1/orders/{order_uuid}`), а не сам UUID —
  иначе кардинальность метрики росла бы с числом заказов.
- **Трейсы** — весь путь заказа складывается в один трейс в Jaeger: HTTP-запрос
  к `order`, его gRPC-вызовы к `inventory` и `payment`, запросы в Postgres, Mongo
  и Redis, а дальше через Kafka — `assembly` и `notification`.
- **Health-проверки** — у каждого сервиса есть healthcheck в compose: HTTP-сервисы
  проверяются по `/health`, gRPC — через `grpc-health-probe`, консьюмеры Kafka —
  по `/health` на порту метрик.

### Как контекст едет через Kafka

Продюсер кладёт в заголовки сообщения `traceparent`
(`platform/pkg/kafka/carrier.go` + `producer/producer.go`), а consumer middleware
`Tracing` достаёт его и делает спан продюсера родителем обработки
(`platform/pkg/middleware/kafka/tracing.go`). Без этого sarama отдаёт обработчику
`session.Context()`, в котором спана нет, и каждый сервис начинал бы собственный
трейс — `assembly` и `notification` просто не находились бы по заказу. Middleware
стоит первой в цепочке, чтобы логирование и бизнес-логика попали внутрь спана.
Склейка закреплена тестом `TestTracingContinuesProducerTrace`.

### Что видно в трейсе

Инструментирован каждый слой, через который проходит запрос:

| Слой            | Чем инструментирован                    |
|-----------------|-----------------------------------------|
| HTTP (сервер)   | `otelhttp`                              |
| gRPC (обе стороны) | `otelgrpc` через `StatsHandler`      |
| Kafka           | своя пара продюсер/middleware, см. выше |
| PostgreSQL      | `otelpgx` на пуле `pgx`                 |
| MongoDB         | `otelmongo` как монитор команд драйвера |
| Redis           | `redisotel`                             |

Параметры SQL-запросов в спаны намеренно **не** пишутся
(`otelpgx.WithIncludeQueryParameters` не включён): среди значений бывают пароли и
персональные данные, а трейсы читает кто угодно с доступом к Jaeger.

## Прогресс по курсу

- [x] Неделя 1 — HTTP на chi, gRPC-контракты, buf
- [x] Неделя 2 — слоистая архитектура, unit-тесты
- [x] Неделя 3 — Docker, Docker Compose, PostgreSQL, MongoDB
- [x] Неделя 4 — конфигурация из окружения, DI-контейнер, graceful shutdown
- [x] Неделя 5 — Kafka: продюсеры, консьюмеры, Telegram-уведомления
- [x] Неделя 6 — JWT, Redis, bcrypt
- [x] Неделя 7 — логи (zap), метрики (Prometheus + Grafana), трейсы (OpenTelemetry + Jaeger)
- [x] Неделя 8 — сборка всей системы: семь модулей, шесть сервисов, один `task up-all`
