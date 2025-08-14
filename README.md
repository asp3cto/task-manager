# Task Manager REST API

REST API для управления задачами на Go с использованием гексагональной архитектуры.

## Архитектура

Проект реализует гексагональную архитектуру:
- **Domain** - бизнес-сущности (Task)
- **Ports** - интерфейсы для репозиториев и сервисов
- **Adapters** - реализации интерфейсов (HTTP обработчики, in-memory репозиторий)
- **Core/Service** - бизнес-логика
- **Logger** - асинхронная система логирования с JSON-выводом

## Структура проекта

```
task-manager/
├── cmd/
│   └── main.go                     # Точка входа приложения
├── internal/
│   ├── domain/
│   │   └── task.go                 # Доменная модель Task
│   ├── ports/
│   │   ├── repository.go           # Интерфейс репозитория
│   │   └── service.go              # Интерфейс сервиса
│   ├── adapters/
│   │   ├── http/
│   │   │   ├── handler.go          # HTTP обработчики
│   │   │   └── server.go           # HTTP сервер с graceful shutdown
│   │   └── repository/
│   │       └── memory.go           # In-memory реализация репозитория
│   ├── core/
│   │   └── service/
│   │       └── task.go             # Бизнес-логика
│   └── logger/
│       ├── async.go                # Асинхронный логгер с JSON-форматом
│       └── config.go               # Конфигурация логгера из переменных окружения
├── go.mod
└── README.md
```

## API Endpoints

### GET /tasks
Получить список всех задач с опциональной фильтрацией по статусу.

**Query Parameters:**
- `status` (optional) - фильтр по статусу: `pending`, `in_progress`, `completed`, `cancelled`

**Пример запроса:**
```bash
curl http://localhost:8080/tasks
curl http://localhost:8080/tasks?status=pending
```

**Пример ответа:**
```json
[
    {
        "id": "1a2b3c4d5e6f7g8h",
        "title": "Выполнить задачу",
        "description": "Описание задачи",
        "status": "pending",
        "created_at": "2023-12-01T10:00:00Z",
        "updated_at": "2023-12-01T10:00:00Z"
    }
]
```

### GET /tasks/{id}
Получить задачу по ID.

**Пример запроса:**
```bash
curl http://localhost:8080/tasks/1a2b3c4d5e6f7g8h
```

**Пример ответа:**
```json
{
    "id": "1a2b3c4d5e6f7g8h",
    "title": "Выполнить задачу",
    "description": "Описание задачи",
    "status": "pending",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
}
```

### POST /tasks
Создать новую задачу.

**Request Body:**
```json
{
    "title": "Название задачи",
    "description": "Описание задачи"
}
```

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Новая задача",
    "description": "Описание новой задачи"
  }'
```

**Пример ответа:**
```json
{
    "id": "1a2b3c4d5e6f7g8h",
    "title": "Новая задача",
    "description": "Описание новой задачи",
    "status": "pending",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
}
```

## Статусы задач

- `pending` - ожидает выполнения
- `in_progress` - в процессе выполнения
- `completed` - завершена
- `cancelled` - отменена

## Логирование

Приложение использует асинхронную систему логирования с JSON-форматом вывода.

### Конфигурация через переменные окружения

- `LOG_LEVEL` - уровень логирования (DEBUG, INFO, WARN, ERROR). По умолчанию: INFO
- `LOG_BUFFER_SIZE` - размер буфера для очереди логов. По умолчанию: 100

### Пример логов
```json
{"time":"2023-12-01T10:00:00Z","level":"INFO","message":"server starting","addr":":8080"}
{"time":"2023-12-01T10:00:05Z","level":"DEBUG","message":"task created successfully","task_id":"1a2b3c4d5e6f7g8h","title":"New Task"}
{"time":"2023-12-01T10:00:10Z","level":"WARN","message":"invalid status parameter","status":"invalid"}
{"time":"2023-12-01T10:00:15Z","level":"ERROR","message":"failed to create task","error":"database connection failed"}
```

## Сборка и запуск

### Требования
- Go 1.24 или выше

### Сборка
```bash
go build -o task-manager cmd/main.go
```

### Запуск
```bash
# Запуск с настройками по умолчанию
./task-manager

# Или запуск из исходного кода
go run cmd/main.go

# Запуск с кастомными настройками
ADDR=:3000 LOG_LEVEL=DEBUG LOG_BUFFER_SIZE=200 ./task-manager
```

### Переменные окружения
- `ADDR` - адрес и порт для прослушивания (по умолчанию: `:8080`)
- `LOG_LEVEL` - уровень логирования: DEBUG, INFO, WARN, ERROR (по умолчанию: `INFO`)
- `LOG_BUFFER_SIZE` - размер буфера логов (по умолчанию: `100`)

### Graceful Shutdown
Сервер поддерживает graceful shutdown. Для остановки используйте Ctrl+C (SIGINT) или отправьте SIGTERM. При завершении все оставшиеся логи будут записаны.

## Примеры использования

### Создание задачи
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Изучить Go",
    "description": "Изучить основы языка Go и создать простое API"
  }'
```

### Получение всех задач
```bash
curl http://localhost:8080/tasks
```

### Получение задач по статусу
```bash
curl http://localhost:8080/tasks?status=pending
```

### Получение задачи по ID
```bash
curl http://localhost:8080/tasks/{task_id}
```

## Ошибки

API возвращает ошибки в формате JSON:
```json
{
    "error": "описание ошибки"
}
```

HTTP статус коды:
- `200` - успешный запрос
- `201` - успешное создание
- `400` - некорректный запрос
- `404` - ресурс не найден
- `405` - метод не разрешен
- `500` - внутренняя ошибка сервера
