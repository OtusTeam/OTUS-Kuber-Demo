# PLG Stack Demo (Promtail + Loki + Grafana)

## 🎯 Цель демо

Развертывание и настройка централизованного логирования в Kubernetes с использованием PLG стека.

## 📦 Компоненты

- **Promtail** - Агент сбора логов
- **Loki** - Система хранения и индексации логов
- **Grafana** - Визуализация и дашборды

## 🚀 Развертывание

### 1. Подготовка k3s кластера

```bash
cd k3s-setup/
./setup-k3s.sh
```

### 2. Установка PLG стека через Helm

```bash
# Добавляем репозитории Helm
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Устанавливаем Loki
helm install loki grafana/loki-stack -f loki/values.yaml

# Устанавливаем Grafana
helm install grafana grafana/grafana -f grafana/values.yaml
```

### 3. Настройка Promtail

```bash
kubectl apply -f promtail/
```

## 📊 Доступ к Grafana

```bash
# Получаем пароль админа
kubectl get secret grafana -o jsonpath="{.data.admin-password}" | base64 --decode

# Проброс портов
kubectl port-forward svc/grafana 3000:80
```

Grafana будет доступна по адресу: http://localhost:3000

## 🔍 Тестирование

1. Создайте тестовое приложение, которое генерирует логи
2. Проверьте, что логи попадают в Loki
3. Настройте дашборды в Grafana
