
# Authorization Service (Go)

Микросервис для аутентификации и управления токенами с использованием PostgreSQL и JWT.

## 📌 Основные возможности

- Генерация и валидация JWT-токенов   
- Защищённые эндпоинты  
- Интеграция с Swagger  
- Поддержка Docker  

## 🚀 Запуск проекта

### Требования
- Go 1.24+  
- PostgreSQL 15+  
- Docker

### 1. Локальный запуск

```bash
# Установите переменные окружения
cp .env.example .env

# Запустите PostgreSQL (Docker)
docker-compose up -d db

# Запустите приложение
go run cmd/main.go
```

### 2. Запуск через Docker

```bash
docker-compose up --build
```

Приложение будет доступно на (http://localhost:8080)

## 📚 API Endpoints

### 🔓 Публичные эндпоинты

- **Получение токенов**  
  `GET /tokens/:user_guid`  
  access и refresh токены устанавливаются в куки.

- **Обновление токенов**  
  `GET /refresh`  
  Получает токены из кук. Обновляет пару токенов (при смене User-Agent сессия удаляется, при смене IP отправляется webhook).

### 🔒 Защищённые эндпоинты

- **Получение GUID пользователя**  
  `GET /user/guid`  
  Возвращает GUID текущего пользователя.

- **Выход из системы**  
  `POST /auth/logout`  
  Удаляет текущую сессию.

## 🔐 Особенности безопасности

- При изменении User-Agent во время обновления токенов:  
  - Сессия удаляется из БД  
  - Все токены становятся недействительными  

- При изменении IP-адреса:  
  - Отправляется webhook 

- Refresh токены хранятся в БД в хешированном виде.

## 🗄️ Структура БД

Таблица `sessions` для хранения сессий:

```sql
CREATE TABLE sessions (
    session_id SERIAL PRIMARY KEY,
    user_guid UUID NOT NULL,
    user_agent TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```

## 🌐 Swagger Documentation

Документация API доступна после запуска сервера:  
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## ⚙️ Конфигурация (.env)

```ini
# Rest
PORT=:8080

# Postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mydb
DB_USER=admin
DB_PASSWORD=admin
SSL_MODE=disable
DB_POOL_MAX_CONNS=10
DB_POOL_MAX_CONN_LIFETIME=300s
DB_POOL_MAX_CONN_IDLE_TIME=150s

# Token
JWT_SECRET=your_secure_secret_here
LEN_TOKEN=32

# Webhook
WEBHOOK_URL=https://example.com/security-alerts
```

## 🛠️ Технологии

- Go 1.24  
- PostgreSQL  
- JWT  
- Swagger  
- Docker  
- Fiber (web framework)
