# Запуск сканирования

```bash
syft gcr.io/distroless/base-debian12
grype gcr.io/distroless/base-debian12
trivy image --severity HIGH,CRITICAL gcr.io/distroless/base-debian12
trivy image --severity HIGH,CRITICAL bkimminich/juice-shop
```

# Harbor Demo (бинарное хранилище: docker images + helm charts)

Развертывание Harbor в minikube с доступом по HTTPS через ingress и полный
цикл `docker login` / `docker push` с хоста.

## 📦 Компоненты

- **Harbor** (helm chart `goharbor/harbor`) — registry для docker-образов и helm-чартов
- **Ingress** (ingress-nginx из minikube addons) — TLS-терминация для `https://harbor.local`
- **Собственный CA + сертификат** (`certs/`) — чтобы docker-клиент доверял registry
- **NodePort сервис** (`manifests/harbor-ui-nodeport.yaml`) — быстрый доступ к UI без /etc/hosts

Ключевой момент: `docker login` и `docker push` работают **по HTTPS с доверенным CA**
(через `/etc/docker/certs.d/`), а не через `insecure-registries`.

## 🚀 Развертывание

### 1. Подготовка minikube

```bash
minikube start --cpus=2 --memory=4096
minikube addons enable ingress

# Имя harbor.local должно указывать на minikube IP
echo "$(minikube ip) harbor.local" | sudo tee -a /etc/hosts
```

### 2. Установка Harbor

```bash
helm repo add harbor https://helm.goharbor.io
helm repo update

helm upgrade --install harbor harbor/harbor \
    --create-namespace \
    --namespace harbor \
    --values harbor/values.yaml
```

### 3. Ingress + TLS + NodePort UI

```bash
kubectl apply -f manifests/
kubectl wait --namespace harbor --for=condition=Ready pod --all --timeout=600s
```

Сгенерировать при сертификаты:

```bash
bash certs/gen-certs.sh && kubectl apply -f manifests/harbor-tls-secret.yaml
```

### 4. Доступ к UI

```bash
# Через ingress (TLS, то же имя, что и для docker):
#   https://harbor.local  (admin / Harbor12345)

# Или без /etc/hosts через NodePort:
minikube service harbor-ui -n harbor
```

## 🐳 Docker login и push с хоста

### 1. Доверие CA для docker-демона

```bash
sudo mkdir -p /etc/docker/certs.d/harbor.local
sudo cp certs/ca.crt /etc/docker/certs.d/harbor.local/ca.crt
sudo systemctl restart docker   # macOS: restart Docker Desktop / colima
```

### 2. Логин

Перед логином убедитесь, что кластер запущен и ingress работает (см. «Устранение неполадок» ниже):

```bash
kubectl -n ingress-nginx get pods   # ingress-nginx-controller — Running
kubectl -n harbor get pods          # все поды — Running
# 401 = это норма (auth challenge registry, TLS и ingress работают):
curl --cacert certs/ca.crt -i https://harbor.local/v2/
```

```bash
docker login harbor.local -u admin
# Password: Harbor12345 (см. harborAdminPassword в harbor/values.yaml)
```

### 3. Проект и push

Push идёт только в существующий проект — создайте `demo` в UI
(Projects → New Project, публичный для простоты) или через API:

```bash
curl --cacert certs/ca.crt -u admin:Harbor12345 \
  -X POST https://harbor.local/api/v2.0/projects \
  -H "Content-Type: application/json" \
  -d '{"project_name":"demo","metadata":{"public":"true"}}'
```

Пушим образ:

```bash
docker pull busybox:latest
docker tag busybox:latest harbor.local/demo/busybox:latest
docker push harbor.local/demo/busybox:latest
```

Проверка в UI: `https://harbor.local` → проект `demo` → репозиторий `busybox`.

### 4. Pull обратно

```bash
docker rmi harbor.local/demo/busybox:latest busybox:latest
docker pull harbor.local/demo/busybox:latest
```

## 🔁 Pull изнутри кластера (containerd в minikube)

Нодам minikube тоже нужно доверять CA (иначе `ImagePullBackOff`):

```bash
minikube ssh "sudo mkdir -p /etc/containerd/certs.d/harbor.local"
cat certs/ca.crt | minikube ssh "cat > /tmp/ca.crt && sudo mv /tmp/ca.crt /etc/containerd/certs.d/harbor.local/ca.crt && sudo systemctl restart containerd"
```

Пример pod с образом из Harbor:

```bash
kubectl run busybox --image=harbor.local/demo/busybox:latest \
  -n default --restart=Never --rm -it -- sh
```

## 🛠 Устранение неполадок

### `dial tcp 192.168.49.2:443: connect: no route to host`

Minikube **остановлен** (перезагрузка машины, `minikube stop`, сон ноутбука):
IP кластера перестаёт существовать, и ядро отвечает `no route to host`.

```bash
minikube status          # host/kubelet/apiserver должны быть Running
minikube start           # если Stopped — запустить и подождать
kubectl -n ingress-nginx get pods   # дождаться ingress-nginx-controller Running
```

После запуска кластера `docker login` работает без изменений — имя `harbor.local`
и сертификаты не зависят от состояния кластера.

### `dial tcp ...:443: connect: connection refused`

Ingress controller ещё не поднялся (или addon не включён). Проверить:

```bash
minikube addons list | grep ingress   # должно быть enabled
minikube addons enable ingress
kubectl -n ingress-nginx get pods
```

### `x509: certificate signed by unknown authority` (от docker)

Docker-демон не видит CA. Скопировать сертификат и перезапустить демона
(шаг «Доверие CA для docker-демона» выше). После изменений в
`/etc/docker/certs.d/` перезапуск docker **обязателен**.

### `harbor.local` открывается не туда после пересоздания кластера

После `minikube delete && minikube start` IP мог измениться — обновить
запись в `/etc/hosts` (и удалить старую):

```bash
sudo sed -i '/harbor.local/d' /etc/hosts
echo "$(minikube ip) harbor.local" | sudo tee -a /etc/hosts
```

### Поды Harbor в CrashLoopBackOff сразу после возобновления кластера

После `minikube stop/start` (или сна машины) часть подов (jobservice, nginx)
может упасть пару раз — это транзиентное поведение при восстановлении
состояния. Подождите 1–2 минуты: `kubectl -n harbor get pods -w`.

## 📁 Структура

```
demo/
├── README.md                       # этот файл
├── harbor/values.yaml              # values чарта: clusterIP + внешние manifests
├── manifests/
│   ├── harbor-tls-secret.yaml      # TLS-секрет (генерируется certs/gen-certs.sh)
│   ├── harbor-ingress.yaml         # ingress: TLS, proxy-body-size=0 для push
│   └── harbor-ui-nodeport.yaml     # NodePort 30300 для быстрого доступа к UI
└── certs/
    ├── gen-certs.sh                # генерация CA + сертификата + secret-манифеста
    ├── ca.crt / ca.key             # CA (ca.crt копируется в /etc/docker/certs.d)
    └── harbor.local.crt / .key / .csr
```

## ⚠️ Примечания

- Пароли/ключи в репозитории — **демо-значения**. Для реального окружения
  используйте cert-manager, external secrets и сгенерированные пароли.
- `harbor.local` в `/etc/hosts` переживёт пересоздание кластера только если
  minikube IP не изменился; после `minikube delete` обновите запись.
- Секреты `ca.key`, `harbor.local.key` и secret-манифест с ключом не должны
  попадать в публичный git — здесь оставлены для простоты демо.
