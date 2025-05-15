#kind
##クラスタ構築
```
kind create cluster --config kind.yaml --name grafana-cluster
```

##クラスタ削除
```
kind delete cluster --name grafana-cluster
```

## Prometheus構築
### リポジトリ登録
```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

### values取得
```
helm show values prometheus-community/prometheus > kubernetes/helm/prometheus/values.yaml
```

### インストール/更新
```
helm upgrade --install --create-namespace --namespace prometheus prometheus prometheus-community/prometheus -f kubernetes/helm/prometheus/values.yaml
```

### ConfigMap設定
```
k apply -f kubernetes/manifest/prometheus/ -R
```

### PortFoward設定
```
export POD_NAME=$(kubectl get pods --namespace prometheus -l "app.kubernetes.io/name=prometheus,app.kubernetes.io/instance=prometheus" -o jsonpath="{.items[0].metadata.name}")
  kubectl --namespace prometheus port-forward $POD_NAME 9090
```

##  Grafana構築
### リポジトリ登録
```
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

### values取得
```
helm show values grafana/grafana > kubernetes/helm/grafana/values.yaml
```

### インストール/更新
```
helm upgrade --install --create-namespace --namespace grafana grafana grafana/grafana -f kubernetes/helm/grafana/values.yaml
```

### manifest追加
```
k apply -f kubernetes/manifest/grafana/ -R
```

### password確認
```
kubectl get secret --namespace grafana grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

### PordForward設定
```
export POD_NAME=$(kubectl get pods --namespace grafana -l "app.kubernetes.io/name=grafana,app.kubernetes.io/instance=grafana" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace grafana port-forward $POD_NAME 3000
```


## FluentBitインストール
### リポジトリ登録
```
helm repo add fluent https://fluent.github.io/helm-charts
helm repo update
```

### values取得
```
helm show values fluent/fluent-bit > kubernetes/helm/fluentbit/values.yaml
```

###インストール
```
helm upgrade --install fluent-bit fluent/fluent-bit --create-namespace --namespace fluentbit -f kubernetes/helm/fluentbit/values.yaml
```


## Lokiインストール
## Values取得
```
helm show values grafana/loki > kubernetes/helm/loki/values.yaml
```

###インストール
```
helm upgrade --install loki grafana/loki --create-namespace --namespace loki -f kubernetes/helm/loki/values.yaml
```

## Alloyインストール
## Values取得
```
helm show values grafana/alloy > kubernetes/helm/alloy/values.yaml
```

###インストール
```
helm upgrade --install alloy grafana/alloy --create-namespace --namespace alloy -f kubernetes/helm/alloy/values.yaml
```