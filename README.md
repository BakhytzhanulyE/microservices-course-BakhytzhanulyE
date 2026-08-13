# microservices-course-BakhytzhanulyE

[![CI](https://github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/actions/workflows/ci.yml/badge.svg)](https://github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/actions/workflows/ci.yml)

Здесь я по шагам собираю систему из нескольких Go-микросервисов: gRPC-взаимодействие, слоистая архитектура, тесты, работа с БД, конфигурация и запуск в Docker.

## Стек

- **Go 1.26+**
- **gRPC** / **protobuf** — общение между сервисами
- **PostgreSQL**, **MongoDB** — хранилища
- **Docker**, **Docker Compose** — локальный запуск
- **golangci-lint** — статический анализ (конфиг в `.golangci.yml`)
- **Task** — запуск команд проекта
- **GitHub Actions** — CI

## Структура

Каждый сервис — отдельный Go-модуль со своим `go.mod`.

| Модуль | Назначение |
|---|---|
| `inventory` | каталог деталей |
| `order` | оформление заказов |
| `payment` | оплата |
| `assembly` | сборка |
| `iam` | аутентификация и доступы |
| `notification` | уведомления |
| `platform` | общий код: логгер, конфиг, вспомогательное |

## Требования

- [Go](https://go.dev/dl/) 1.26 или новее
- [Task](https://taskfile.dev/installation/) — `brew install go-task`

Линтер и форматтеры ставить руками не нужно: они скачиваются в `./bin` нужных версий командами ниже.

## Команды

```bash
task install-formatters      # поставить gofumpt и gci в ./bin
task format                  # отформатировать код и импорты
task install-golangci-lint   # поставить линтер в ./bin
task lint                    # проверить код линтером
```

Версии инструментов зафиксированы в `Taskfile.yml` — локально и в CI используются одинаковые.

## CI/CD

GitHub Actions запускается на каждый push и pull request:

- **`.github/workflows/ci.yml`** — линтинг и проверка безопасности
- **`.github/workflows/lint-reusable.yml`** — переиспользуемый workflow для линтинга
- **`.github/scripts/extract-versions.sh`** — тянет версии инструментов из `Taskfile.yml`, чтобы CI и локальная машина не разъезжались

## Прогресс

- [ ] Неделя 1 — HTTP и gRPC
- [ ] Неделя 2 — слои и unit-тесты
- [ ] Неделя 3 — Docker, PostgreSQL, MongoDB
- [ ] Неделя 4 — конфигурация и DI
- [ ] Неделя 5
- [ ] Неделя 6
- [ ] Неделя 7
- [ ] Неделя 8
