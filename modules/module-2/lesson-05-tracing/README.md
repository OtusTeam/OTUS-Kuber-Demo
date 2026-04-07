# Demo zerocode auto-intrumentation

## Dependencies

1. k8s cluster, например k3s, minikube или облако

2. [Cert manager](https://cert-manager.io/)

3. [Ingress-nginx](https://kubernetes.github.io/ingress-nginx/deploy/)

в configmaps `ingress-nginx-controller` сделать 2 настройки

```yaml
data:
  enable-opentelemetry: "true"
  otlp-collector-host: demo-collector.default
```

4. [OpenTelemtry Operator](https://opentelemetry.io/docs/platforms/kubernetes/operator/#getting-started)

## Установка

1. Установить Jaeger любым способом, в примере через helm

```bash
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm install jaeger jaegertracing/jaeger
```

2. Применить объекты k8s opentelementry operator

```bash
kubectl apply -f instrumetation.yaml
kubectl apply -f collector.yaml
```

3. Применить объекты приложения

```bash
kubectl apply -f rbac.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
```

## Demo

Сделать перенаправление портов на лоальный хост для ingress

```bash
kubectl port-forward --namespace=ingress-nginx service/ingress-nginx-controller 8080:80
```

и cделать запрос в ingress

```bash
curl --resolve demo.localdev.me:8080:127.0.0.1 http://demo.localdev.me:8080
```

Сделать перенаправление портов на лоальный хост для jaeger, чтобы посмотреть trace

```bash
kubectl port-forward --namespace default service/jaeger 16686:16686
```

## Дальнейшие настройки

* [Сollector](https://opentelemetry.io/docs/collector/configuration/)
* [Python](https://opentelemetry.io/docs/zero-code/python/)
* [Igress-nginx](https://kubernetes.github.io/ingress-nginx/user-guide/third-party-addons/opentelemetry/)

## Удаление

```bash
kubectl delete all --all
kubectl delete -f instrumetation.yaml
kubectl delete -f collector.yaml
```
