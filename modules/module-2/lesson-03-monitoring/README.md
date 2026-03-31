# Prometheus

## Установка

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

```bash
kubectl create namespace monitoring
```

```
helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring
```

```
helm install prometheus oci://ghcr.io/prometheus-community/charts/kube-prometheus-stack --namespace monitoring
```

```
kubectl port-forward svc/prometheus-kube-prome-prometheus 9090:9090 --namespace monitoring
```

```
kubectl get secret --namespace monitoring prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 --decode
```

```
kubectl port-forward svc/prometheus-grafana 3000:3000 --namespace monitoring
```



## Конфигурирование

````bash
helm show values prometheus-community/kube-prometheus-stack > values.yaml
````

Можно поменять необходимые опции, и запустить (например, отключить grafana)

```bash
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring -f values.yaml
```



## Promql

Количество запросов в период времени 1m

```
round(sum by (code) (increase(apiserver_request_total{job="apiserver"}[1m])))
```



## Примеры alterts

https://samber.github.io/awesome-prometheus-alerts/rules#postgresql

