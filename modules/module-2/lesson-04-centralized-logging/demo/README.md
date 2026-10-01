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
helm repo add grafana-community https://grafana-community.github.io/helm-charts
helm repo update

# Устанавливаем Loki и Grafana
helm upgrade --install \
    loki grafana-community/loki  \
    --create-namespace \
    --namespace loki \
    --values loki/values.yaml

# Устанавливаем Grafana с преднастроенным Loki datasource
helm upgrade --install \
    grafana grafana-community/grafana \
    --namespace loki \
    --values grafana/values.yaml

# Promtail удалён из актуального индекса grafana-community (deprecated в пользу Alloy),
# финальная версия чарта доступна через legacy-репозиторий
helm repo add grafana-legacy https://grafana.github.io/helm-charts

# Устанавливаем Promtail (DaemonSet, собирает логи всех подов)
helm upgrade --install \
    promtail grafana-legacy/promtail \
    --namespace loki \
    --values promtail/values.yaml
```

> ⚠️ **Важно для minikube (docker driver):** у всех «нод» кластера одно общее ядро хоста,
> а лимит `fs.inotify.max_user_instances` по умолчанию 128. Три пода promtail суммарно
> открывают больше inotify-watcher'ов и падают с `too many open files`.
> Поднимите лимит перед установкой promtail:
>
> ```bash
> minikube ssh "sudo sysctl -w fs.inotify.max_user_instances=1024"
> ```
>
> (значение не переживает перезапуск minikube — при пересоздании кластера повторить)

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

Grafana разворачивается с уже настроенным Loki datasource (`grafana/values.yaml`).

```bash
# Открываем UI (minikube сам откроет браузер)
minikube service grafana -n loki

# или напрямую: http://<minikube ip>:30300
minikube ip

# или через port-forward
kubectl port-forward svc/grafana 3000:80 -n loki
```

Логин/пароль: `admin` / `admin` (демо-значения заданы в `grafana/values.yaml`).

## 🔍 Просмотр логов Loki в Grafana

1. Откройте **Explore** → выберите datasource **Loki**
2. Label browser: `{job="test"}` — тестовые логи, отправленные вручную через API
3. Через **LogQL** можно фильтровать: `{job="test"} |= "hello"`
4. Логи, собранные **promtail** (весь кластер): селекторы по меткам
   - `{pod="loki-0"}` — логи самого Loki
   - `{namespace="kube-system", container="coredns"}` — логи CoreDNS
   - `{app="grafana"}` — логи Grafana
   - метки: `namespace`, `pod`, `container`, `app`, `job`, `node_name`

Отправить тестовые логи в Loki:

```bash
kubectl run loki-push-test --image=curlimages/curl --restart=Never --rm -it --command -- \
  curl -s -H "Content-Type: application/json" -XPOST \
  http://loki-gateway.loki.svc.cluster.local/loki/api/v1/push \
  --data-raw "{\"streams\": [{\"stream\": {\"job\": \"test\"}, \"values\": [[\"$(date +%s)000000000\", \"hello from PLG demo\"]]}]}"
```

## 🔍 Тестирование

1. Создайте тестовое приложение, которое генерирует логи
2. Проверьте, что логи попадают в Loki
3. Настройте дашборды в Grafana
