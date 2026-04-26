# Grafana dashboard

Файл `game-catalog-dashboard.json` импортируется в Grafana Cloud через:

`Dashboards` -> `New` -> `Import` -> `Upload dashboard JSON file`.

При импорте выбери Prometheus datasource, в который пишет облачный Prometheus/Grafana Agent. Dashboard ожидает namespace `game-catalog`.

Перед импортом или обновлением dashboard примени Kubernetes-манифесты, чтобы Prometheus Operator начал явно собирать backend `/metrics`:

```bash
kubectl apply -k k8s
kubectl get podmonitor backend -n game-catalog
```

## Что показывает dashboard

- состояние pod'ов и количество backend replicas;
- CPU и memory по pod'ам из kubelet cAdvisor metrics;
- network I/O по node metrics;
- HTTP request rate по backend pod'ам;
- таблицу запросов приложения по `pod`, `method`, `path`, `status`;
- p95 latency по backend pod'ам;
- 5xx error rate;
- текущие in-flight requests;
- сравнение backend CPU usage/request/limit.

## Базовые проверки в Explore

Если инфраструктурные панели показывают `No data`, проверь, что в Prometheus есть метрики cAdvisor с такими labels:

```promql
sum by (pod) (
  rate(container_cpu_usage_seconds_total{
    job="kubelet",
    metrics_path="/metrics/cadvisor",
    namespace="game-catalog"
  }[5m])
)
```

В текущем kube-prometheus-stack cAdvisor может не отдавать `container_network_receive_bytes_total`/`container_network_transmit_bytes_total`, поэтому dashboard использует node-level network metrics:

```promql
sum by (instance) (
  rate(node_network_receive_bytes_total{device!="lo"}[5m])
)
```

Если панели приложения показывают `No data`, проверь scrape backend:

```promql
sum by (pod, method, path, status) (
  rate(game_catalog_http_requests_total{namespace="game-catalog"}[5m])
)
```

Если запрос выше не возвращает `pod`, значит Prometheus не добавляет Kubernetes target labels к annotation scrape. В этом случае в dashboard нужно заменить фильтр `pod=~"backend-.*"` на label, который реально есть в Explore, например `instance`, `service`, `job` или `k8s_pod_name`.
