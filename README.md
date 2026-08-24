# MoodFly - авиационная система управления

## Описание

REST API сервис для управления авиационной деятельностью (Go, net/http, Docker, Redis, JWT)

## Технологический стек

- Golang
- net/http + ServeMux
- PostgreSQL
- Redis
- JWT
- Taskfile
- Docker


## Установка и запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/GERKULE5/MoodFly.git && cd MoodFly
```

### 2. Настройка окружения

Переименуйте .env.example -> .env

### 3. Установка зависимостей

```bash
go mod download
```

### Запуск

```bash
task dev
```

### Ссылки

Сервис доступен по адресу: http://localhost:8080


## API

### Ping

- GET ping/

### Users

- POST /users
- GET /users
- GET /users/{id}
- PUT /users/{id}
- DELETE /users/{id}

### Auth

- POST /register
- POST /login 
- POST /refresh
- POST /logout

### Aircrafts

- POST /aircrafts
- GET /aircrafts
- GET /aircrafts/{id}
- PUT /aircrafts/{id}
- DELETE /aircrafts/{id}

### Flights

- POST /flights
- GET /flights
- GET /flights/{id}
- PUT /flights/{id}
- DELETE /flights/{id}
