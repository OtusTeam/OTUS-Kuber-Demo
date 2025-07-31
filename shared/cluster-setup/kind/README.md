# Kind Cluster Configurations

Конфигурационные файлы для создания кластеров Kind для демонстраций OTUS.

## Quick Start

### Установка kind

```bash
# macOS
brew install kind

# Linux
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

### Создание простого кластера

```bash
kind create cluster --name otus-demo
```

### Создание кластера с конфигурацией

```bash
# Одноузловой кластер
kind create cluster --name otus-demo --config single-node.yaml

# Многоузловой кластер
kind create cluster --name otus-demo --config multi-node.yaml

# Кластер с Ingress
kind create cluster --name otus-demo --config ingress.yaml

# Кластер с Load Balancer
kind create cluster --name otus-demo --config loadbalancer.yaml

# HA кластер с несколькими control-plane
kind create cluster --name otus-demo-ha --config ha-cluster.yaml

# Кластер для разработки с дополнительными настройками
kind create cluster --name otus-demo-dev --config development.yaml
```

## Доступные конфигурации

### single-node.yaml
Простой одноузловой кластер для базовых демонстраций.

### multi-node.yaml
Кластер с 1 control plane и 2 worker узлами для демонстрации распределенных приложений.

### ingress.yaml
Кластер с настройкой портов для Ingress Controller (порты 80 и 443).

### loadbalancer.yaml
Кластер с поддержкой LoadBalancer сервисов (можно использовать с MetalLB).

### ha-cluster.yaml
Высокодоступный кластер с 3 control plane и 3 worker узлами.

### development.yaml
Кластер для разработки с дополнительными портами и монтированием локальных директорий.

## Управление кластером

```bash
# Список кластеров
kind get clusters

# Получение kubeconfig
kind get kubeconfig --name otus-demo

# Удаление кластера
kind delete cluster --name otus-demo

# Загрузка образа в кластер
kind load docker-image my-image:tag --name otus-demo
```

## Полезные команды

```bash
# Подключение к узлу
docker exec -it otus-demo-control-plane bash

# Проверка логов
docker logs otus-demo-control-plane

# Список контейнеров кластера
docker ps --filter "label=io.x-k8s.kind.cluster=otus-demo"
```

## Ссылки

- [Kind Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Kind Configuration](https://kind.sigs.k8s.io/docs/user/configuration/)
- [Ingress with Kind](https://kind.sigs.k8s.io/docs/user/ingress/)
- [KIND github with examples](https://github.com/kubernetes-sigs/kind)
