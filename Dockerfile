# ==============================
# 1) BUILD STAGE
# ==============================
FROM golang:1.23-alpine AS builder

WORKDIR /service

# 1.1. Скачиваем модули
COPY go.mod go.sum ./
RUN go mod download

# 1.2. Копируем весь код сразу — embed.FS подтянет папку internal/assets/migrations
COPY . .

# 1.3. Собираем главный бинарь, в который уже вшиты миграции
RUN CGO_ENABLED=0 GOOS=linux go build \
    -o main \
    ./cmd/chains-auth

# ==============================
# 2) FINAL STAGE
# ==============================
FROM alpine:latest

WORKDIR /service

# 2.1. Минимальные утилиты (для логов, TLS и т.п.)
RUN apk add --no-cache ca-certificates

# 2.2. Копируем только Go-бинарь и конфиг
COPY --from=builder /service/main .
COPY config_docker.yaml .

# 2.3. Переменные окружения
ENV KV_VIPER_FILE=/service/config_docker.yaml
EXPOSE 8001

# 2.4. Точка входа — запускаем ваш сервис с обычным HTTP-сервером
CMD ["./main", "run", "service"]
