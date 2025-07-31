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
k3d cluster create --config multi-node.yaml
```

### 2. Установка PLG стека через Helm

https://grafana.com/docs/loki/latest/setup/install/helm/install-monolithic/

```bash
# Добавляем репозитории Helm
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Устанавливаем Loki и Grafana
helm upgrade --install \
    loki grafana/loki \
    --create-namespace \
    --namespace logs \
    --values loki/values.yaml
```

* Get loki password
    ```bash
    passw=$(kubectl get secret --namespace logs loki-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo)
    printf "L: admin\nP: %s\n" "$passw"
    ```


* Share traffic to test Loki 3100 port
    ```bash
    kubectl port-forward svc/loki-gateway 3100:3100 -n logs
    ```

* Verify that Loki did receive the data using the following command:
    ```bash
    curl -H "Content-Type: application/json" -XPOST -s "http://127.0.0.1:3100/loki/api/v1/push"  \
    --data-raw "{\"streams\": [{\"stream\": {\"job\": \"test\"}, \"values\": [[\"$(date +%s)000000000\", \"fizzbuzz\"]]}]}"
    ```

* More complex log message
    ```bash
    url="https://10b1-94-19-17-241.ngrok-free.app"
    curl -H "Content-Type: application/json" -XPOST -s "$url/loki/api/v1/push"  \
    --data-raw "{\"streams\": [{\"stream\": {\"job\": \"external\"}, \"values\": [[\"$(date +%s)000000000\", \"$(whoami);$(pwd) logs\"]]}]}"
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
