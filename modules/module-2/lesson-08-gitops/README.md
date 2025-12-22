## Demo GitOps и инструменты поставки



## Flux

### Dependencies



### Install

установить cli утилиту

````bash
curl -s https://fluxcd.io/install.sh | sudo bash
````

создать репозиторий в gitlab

```bash
flux bootstrap gitlab --hostname=gitlab.com --owner=w1ndblow --repository=flux --branch=main   --path=clusters/minikube --token-auth --network-policy=false --personal
```

Если команда выполнилась удачно должны быть сгенерированы `gotk-components.yaml`, `gotk-sync.yaml`, `kustomization.yaml` в директории `clusters/minikube`

Добавить файлы из этого репозитория из директории `flux` по нужным путям, чтобы получилась следующая структура

```bash
.
├── apps
│   ├── deployment.yaml
│   ├── kustomization.yaml
│   ├── namespace.yaml
│   └── service.yaml
└── clusters
    └── minikube
        ├── flux-system
        │   ├── gotk-components.yaml
        │   ├── gotk-sync.yaml
        │   ├── kustomization.yaml
        │   └── ui.yaml
        └── kustomizaton.yaml
```

Оправить изменения в gitlab и смотреть за результатом
