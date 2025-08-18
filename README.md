# OTUS Kubernetes Course Demo Repository

[![Course](https://img.shields.io/badge/course-OTUS_Kubernetes-orange.svg)](https://otus.ru/lessons/infrastrukturnaya-platforma-na-osnove-kubernetes/)
[![Pre-commit](https://github.com/OtusTeam/OTUS-Kuber-Demo/actions/workflows/pre-commit.yml/badge.svg)](https://github.com/OtusTeam/OTUS-Kuber-Demo/actions/workflows/pre-commit.yml)

Демонстрационные материалы для курса **"Инфраструктурная платформа на основе Kubernetes"** от OTUS.

## 🏗️ Структура репозитория

```
.
├── docs/           # Документация курса
├── examples/       # Общие примеры
├── assets/         # Медиа файлы
├── shared/         # Общие конфигурации кластеров
├── modules/        # Модули курса
└── README.md       # Этот файл
```

## 📚 Модули курса

| Модуль | Название | Уроки |
|--------|----------|-------|
| 1 | Основы Kubernetes | 1-3 |
| 2 | Экосистема Kubernetes | 4-14 |
| 3 | Безопасность и мониторинг | 15-19 |
| 4 | Продвинутые темы | 20-23 |

## 🚀 Быстрый старт

1. Клонируйте репозиторий:
   ```bash
   git clone <repository-url>
   cd OTUS-Kuber-Demo
   ```

2. Перейдите к нужному модулю:
   ```bash
   cd modules/module-2/lesson-04-centralized-logging
   ```

3. Следуйте инструкциям в README модуля

## 🔧 Настройка кластера

Для демонстраций используйте готовые конфигурации кластеров из директории `shared/`:

```bash
# Создание Kind кластера для демо
cd shared/cluster-setup/kind/
kind create cluster --name otus-demo --config ingress.yaml

# Возврат к уроку
cd ../../../modules/module-2/lesson-04-centralized-logging/
```

Доступные типы кластеров:
- **Kind** - быстрое развертывание в Docker контейнерах
- **k3s** - легковесный Kubernetes для локальной разработки
- **Minikube** - классический инструмент для локальной разработки
- **Yandex Cloud** - managed Kubernetes в облаке

## 🛠️ Требования

- Docker Desktop или Podman
- kubectl
- helm 3.x
- k3/minkube/kind/etc.. (для демонстраций)

## 📖 Документация

Подробная документация доступна в директории [`docs/`](./docs/).

Конфигурации кластеров и инструкции по их настройке доступны в [`shared/`](./shared/).

## 🤝 Вклад в проект

Читайте [CONTRIBUTING.md](./docs/CONTRIBUTING.md) для получения информации о том, как внести вклад в проект.

## 📞 Поддержка

- 📧 Email: support@otus.ru
- 🌐 Website: https://otus.ru
- 💬 Telegram: [OTUS. Сопровождение преподавателей](https://t.me/c/1290549066/8684)

---
**OTUS** | Практические курсы для IT-специалистов
