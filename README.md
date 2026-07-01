# Subscriptions Service

[![Coverage](https://img.shields.io/badge/coverage-79.4%25-brightgreen)](./coverage.out)

REST-сервис для учета онлайн-подписок пользователей и подсчета суммарной стоимости подписок за выбранный период.

## Стек

Go 1.25, PostgreSQL 18, pgx, net/http, Docker Compose, OpenAPI 3.0.

## Архитектура

Проект разделен на слои:

- `domain` — сущность подписки, месяц как объект-значение, доменные ошибки и правила;
- `application` — сценарии создания, обновления, удаления, получения списка и расчета стоимости;
- `repository` — интерфейс хранилища находится на стороне application-слоя;
- `infrastructure` — PostgreSQL-реализация репозитория, пул соединений, логирование;
- `transport` — HTTP-обработчики, маршрутизация, middleware;
- `config` — загрузка настроек из YAML и переменных окружения `.env`;
- `migrations` — SQL-миграции;
- `docs` — OpenAPI-спецификация.

## Запуск

```bash
docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`.

Проверка:

```bash
curl http://localhost:8080/health
```

Для локального запуска чувствительные параметры берутся из `.env`, а `docker-compose.yml` подключает этот файл через `env_file`.

## Миграции

При запуске PostgreSQL через Docker Compose миграция из `migrations/001_create_subscriptions.sql` применяется автоматически через `/docker-entrypoint-initdb.d`.

Для промышленного запуска лучше использовать отдельный инструмент миграций: goose, golang-migrate или tern.

## Хранение данных

Данные PostgreSQL сохраняются в именованном volume `subscriptions_pgdata`.

- `docker compose down` останавливает контейнеры, но не удаляет данные;
- `docker compose down -v` или `docker volume rm subscriptions_pgdata` удаляют данные полностью.

## Переменные окружения

- `HTTP_ADDR` — адрес HTTP-сервера, по умолчанию `:8080`;
- `POSTGRES_DSN` — строка подключения к PostgreSQL;
- `POSTGRES_MAX_CONNS` — максимальное число соединений в пуле;
- `POSTGRES_MIN_CONNS` — минимальное число соединений в пуле;
- `POSTGRES_MAX_CONN_LIFETIME` — максимальное время жизни соединения, например `1h`.

Пример локального файла лежит в `.env.example`. Секреты должны храниться только в `.env`.

## API

### `GET /health`
Проверка состояния сервиса.

Ответ:

```json
{ "status": "ok" }
```

### `POST /api/v1/subscriptions`
Создание подписки.

Тело запроса:

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```

### `GET /api/v1/subscriptions`
Список подписок с фильтрами `user_id`, `service_name`, `limit`, `offset`.

### `GET /api/v1/subscriptions/{id}`
Получение подписки по UUID.

### `PUT /api/v1/subscriptions/{id}`
Обновление подписки.

### `DELETE /api/v1/subscriptions/{id}`
Удаление подписки.

### `GET /api/v1/subscriptions/total-cost`
Расчет суммарной стоимости за период.

Параметры:

- `from=07-2025`
- `to=12-2025`
- `user_id` и `service_name` опциональны

## API-примеры

Создание подписки:

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H 'Content-Type: application/json' \
  -d '{"service_name":"Yandex Plus","price":400,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025"}'
```

Получение списка:

```bash
curl 'http://localhost:8080/api/v1/subscriptions?limit=20&offset=0'
```

Подсчет стоимости:

```bash
curl 'http://localhost:8080/api/v1/subscriptions/total-cost?from=07-2025&to=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus'
```

## Тесты

Текущее покрытие: `79.4%`

```bash
go test ./...
```

Покрытие:

```bash
make coverage
```

Или вручную:

```bash
go test -coverpkg=./... -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Smoke-тест на весь HTTP-стек:

```bash
go test ./tests/http -run TestSmokeAllEndpoints -count=1
```

## OpenAPI

Спецификация лежит в `docs/openapi.yaml`.

## Возможные улучшения

- добавить инструмент управляемых миграций вместо init-скрипта PostgreSQL;
- добавить интеграционные тесты repository-слоя через testcontainers;
- добавить трассировку и метрики Prometheus;
- добавить курсорную пагинацию вместо offset для больших таблиц;
- добавить генерацию серверного кода из OpenAPI;
- добавить CI с проверками `go test`, `go vet`, `golangci-lint`;
- добавить отдельные healthcheck-команды для миграций и базы данных.
