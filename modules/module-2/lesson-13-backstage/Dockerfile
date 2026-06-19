# === ЭТАП 1: Сборка всего монорепозитория ===
FROM node:24-trixie-slim

# Устанавливаем зависимости для работы git и сборки нативных модулей
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 sqlite3 libsqlite3-dev g++ make git ca-certificates python3-pip bash gettext-base && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Клонируем репозиторий Backstage
RUN git clone --branch v1.49.5 --depth 1 https://github.com/backstage/backstage.git .


# Активируем Corepack, ставим зависимости, добавляем better-sqlite3 и собираем бэкенд

RUN pip3 install cookiecutter --break-system-packages



# Теперь Yarn видит всю структуру монорепозитория и успешно ставит продакшн-зависимости

COPY app-config.yaml ./

RUN echo "backstage" | npx @backstage/create-app@latest

WORKDIR /app/backstage

# RUN yarn backstage-cli versions:bump

# Копируем конфигурацию для запуска приложения
COPY app-config.yaml ./
COPY configs/backend.index.ts ./packages/backend/src/index.ts
COPY configs/frontend.tsx ./packages/app/src/App.tsx
COPY configs/sidebar.tsx ./packages/app/src/modules/nav/Sidebar.tsx

RUN corepack enable && \
    yarn add sqlite3 && \
    yarn --cwd packages/backend add @backstage/plugin-scaffolder-backend-module-gitlab && \
    yarn --cwd packages/backend add @backstage/plugin-auth-backend-module-guest-provider && \
    yarn --cwd packages/backend add @backstage/plugin-kubernetes-backend && \
    yarn --cwd packages/app add @backstage/plugin-kubernetes@0.12.20-next.1 && \
    yarn install --immutable && \
    yarn tsc && \
    yarn build:all

# Копируем entrypoint скрипт
COPY entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 7007
EXPOSE 3000

# Используем entrypoint для обработки переменных окружения и рендеринга конфига
ENTRYPOINT ["entrypoint.sh"]
CMD []
