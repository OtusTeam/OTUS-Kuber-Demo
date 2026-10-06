#!/bin/bash

docker_compose() {
    # Проверяем, доступна ли команда docker-compose как отдельная утилита
    if command -v docker-compose &> /dev/null; then
        docker-compose "$@"
    else
        docker compose "$@"
    fi
}
