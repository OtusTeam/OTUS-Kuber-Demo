#!/bin/bash

# Функция для вывода сообщений с временной меткой
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

log "Starting Backstage application..."

# Проверяем наличие app-config.yaml шаблона
if [ ! -f "/app/backstage/app-config.yaml" ]; then
    log "ERROR: app-config.yaml not found"
    exit 1
fi

# Создаем временный файл для сгенерированного конфига
TEMP_CONFIG="/tmp/app-config-rendered.yaml"

log "Processing environment variables and rendering configuration..."

# Используем envsubst для подстановки переменных окружения в конфигурационный файл
# envsubst заменяет ${VAR} на значения переменных окружения
envsubst < /app/backstage/app-config.yaml > "$TEMP_CONFIG"

# Проверяем, что файл был создан успешно
if [ ! -f "$TEMP_CONFIG" ]; then
    log "ERROR: Failed to render configuration"
    exit 1
fi

log "Configuration rendered successfully"
log "Starting Backstage with rendered configuration..."

# Запускаем приложение с сгенерированным конфигом
exec yarn start --config "$TEMP_CONFIG"
