# Каталог видеоигр

Учебный fullstack-проект для лабораторной работы по DevOps. Приложение позволяет хранить список видеоигр и выполнять над ними CRUD-операции: создавать, просматривать, редактировать и удалять записи.

## Функциональность

- добавление новой игры
- получение списка игр
- просмотр карточки одной игры
- редактирование игры
- удаление игры
- фильтрация по `status` и `genre`

## Стек технологий

- Backend: Go, Echo, PostgreSQL
- Frontend: React, TypeScript, Vite
- Backend tests: Go `testing`, `httptest`, `gomock`
- Frontend tests: Vitest, React Testing Library
- CI: GitHub Actions

## Структура проекта

```text
.
├── .github/workflows/ci.yml
├── docker-compose.yml
├── backend
│   ├── app.go
│   ├── app_test.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── internal
│   │   ├── game
│   │   ├── httpapi
│   │   ├── media
│   │   └── storage/postgres
│   │       └── migrations
│   ├── main.go
│   └── go.sum
├── frontend
│   ├── Dockerfile
│   ├── index.html
│   ├── nginx
│   │   └── default.conf
│   ├── package.json
│   ├── tsconfig.json
│   ├── src
│   │   ├── api.ts
│   │   ├── App.test.tsx
│   │   ├── App.tsx
│   │   ├── components
│   │   │   ├── GameDetails.tsx
│   │   │   ├── GameForm.tsx
│   │   │   └── GameList.tsx
│   │   ├── constants.ts
│   │   ├── index.css
│   │   ├── main.tsx
│   │   ├── test/setup.ts
│   │   └── types.ts
│   └── vite.config.ts
└── README.md
```

## Модель данных

Сущность `Game`:

- `id` — идентификатор
- `title` — название игры
- `genre` — жанр
- `platform` — платформа
- `releaseYear` — год выпуска
- `rating` — оценка от 1 до 10
- `status` — `planned`, `playing`, `completed`
- `description` — краткое описание игры
- `imagePath` — путь к загруженной обложке внутри `/blob/...`

## Запуск backend

Требование: Go 1.26.1 или новее.

Требование по БД: запущенный PostgreSQL. По умолчанию backend использует строку подключения:

```text
postgres://postgres:postgres@localhost:5432/game_catalog?sslmode=disable
```

Если нужно, можно передать свою строку через переменную окружения `DATABASE_URL`.

Для локального запуска PostgreSQL можно поднять через Docker Compose из корня проекта:

```bash
docker compose up -d postgres
```

```bash
cd backend
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/game_catalog?sslmode=disable"
go mod tidy
go run .
```

Backend стартует на `http://localhost:3000`.

При первом запуске:

- автоматически применяются SQL-миграции из `backend/internal/storage/postgres/migrations`
- автоматически добавляются несколько демонстрационных игр
- автоматически используется папка `backend/blob` для обложек игр

Доступные команды backend:

```bash
go run .
go build ./...
go test ./...
```

## Запуск через Docker Compose

Чтобы поднять сразу PostgreSQL, backend и frontend:

```bash
docker compose up --build
```

После старта сервисы будут доступны по адресам:

- frontend: `http://localhost:5173`
- backend: `http://localhost:3000`
- postgres: `localhost:5432`

Внутри Docker Compose frontend проксирует запросы `/api/*` и `/blob/*` в сервис `backend`, поэтому отдельный `VITE_API_URL` для контейнерного запуска не нужен.

## Запуск в minikube

Требования: установлены `minikube`, `kubectl` и Docker.

1. Запустить minikube:

```bash
minikube start -p minikube
kubectl config use-context minikube
```

Если context `minikube` не появился, обновить его можно командой:

```bash
minikube update-context -p minikube
kubectl config use-context minikube
```

2. Собрать образы backend и frontend локальным Docker:

```bash
docker build -t game-catalog-backend:local ./backend
docker build -t game-catalog-frontend:local ./frontend
```

3. Загрузить локальные образы в minikube:

```bash
minikube image load game-catalog-backend:local
minikube image load game-catalog-frontend:local
```

4. Применить Kubernetes-манифесты:

```bash
kubectl apply -k k8s
```

5. Дождаться запуска всех deployment:

```bash
kubectl rollout status deployment/postgres -n game-catalog
kubectl rollout status deployment/backend -n game-catalog
kubectl rollout status deployment/frontend -n game-catalog
```

6. Открыть frontend:

```bash
minikube service frontend -n game-catalog
```

В Kubernetes frontend, как и в Docker Compose, проксирует `/api/*` и `/blob/*` на внутренний сервис `backend:3000`, поэтому отдельный внешний backend URL не нужен.

### Горизонтальное масштабирование backend

Для backend настроен `HorizontalPodAutoscaler`: Kubernetes держит минимум 1 pod и может поднять до 5 pod'ов, если средняя нагрузка CPU станет выше 15% от `resources.requests.cpu`.

В `backend` задано:

```yaml
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 256Mi
```

Цель HPA `averageUtilization: 15` означает: если backend в среднем потребляет больше 15% от `100m`, то есть больше примерно `15m` CPU на pod, Kubernetes начнет добавлять pod'ы.

Для minikube нужно включить Metrics Server:

```bash
minikube addons enable metrics-server
kubectl top pods -n game-catalog
```

Проверка HPA:

```bash
kubectl get hpa -n game-catalog
kubectl describe hpa backend -n game-catalog
```

Для нагрузочного тестирования backend можно использовать `hey`. Установить:

```bash
brew install hey
```

В отдельном терминале пробросить backend наружу для ручной проверки:

```bash
kubectl port-forward service/backend 3000:3000 -n game-catalog
```

Локальная нагрузка через `port-forward` подходит для быстрой проверки endpoint'а, но не является лучшим способом проверки HPA: port-forward к Service может работать через один выбранный pod, поэтому новые pod'ы не обязательно начнут получать этот же поток запросов.

Для быстрой ручной проверки endpoint'а:

```bash
hey -z 2m -c 20 'http://localhost:3000/api/load/cpu?duration=500ms'
```

Для проверки HPA лучше запускать нагрузку внутри кластера на Kubernetes Service:

```bash
kubectl run backend-load-test -n game-catalog --rm -i --restart=Never \
  --image=rakyll/hey -- \
  -z 5m -c 50 'http://backend:3000/api/load/cpu?duration=500ms'
```

В другом терминале наблюдать масштабирование:

```bash
kubectl get hpa -n game-catalog -w
kubectl get pods -n game-catalog -w
```

### Метрики Prometheus и Grafana Cloud

Backend отдает Prometheus-метрики на:

```text
/metrics
```

Проверить локально внутри кластера:

```bash
kubectl port-forward service/backend 3000:3000 -n game-catalog
curl http://localhost:3000/metrics
```

Основные метрики приложения:

- `game_catalog_http_requests_total` — количество HTTP-запросов с labels `method`, `path`, `status`
- `game_catalog_http_request_duration_seconds` — длительность HTTP-запросов
- `game_catalog_http_requests_in_flight` — текущие активные HTTP-запросы
- стандартные Go/process metrics: `go_*`, `process_*`

В pod template backend добавлены scrape-аннотации:

```yaml
k8s.grafana.com/scrape: "true"
k8s.grafana.com/job: game-catalog-backend
k8s.grafana.com/metrics.path: /metrics
k8s.grafana.com/metrics.portNumber: "3000"
```

Для Grafana Cloud удобнее всего установить Grafana Kubernetes Monitoring Helm chart через мастер настройки в Grafana Cloud. В нем нужно включить:

- Kubernetes infrastructure metrics
- Annotation autodiscovery для application metrics
- Pod logs, если нужно видеть логи по pod'ам рядом с метриками

После подключения в Grafana Cloud можно смотреть:

- CPU/memory по pod'ам backend
- количество pod'ов deployment/backend
- HTTP-запросы по `pod`, `method`, `path`, `status`
- latency по `game_catalog_http_request_duration_seconds`

Полезные команды для диагностики:

```bash
kubectl get pods,svc,pvc -n game-catalog
kubectl logs deployment/backend -n game-catalog
kubectl logs deployment/frontend -n game-catalog
kubectl logs deployment/postgres -n game-catalog
```

Удалить приложение из minikube:

```bash
kubectl delete namespace game-catalog
```

## Запуск frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend стартует на `http://localhost:5173`.

Во время разработки запросы `/api/*` проксируются в backend на `http://localhost:3000`. Frontend написан на TypeScript.

Доступные команды frontend:

```bash
npm run dev
npm run build
npm test
```

## Запуск тестов

Backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd frontend
npm test
```

## CI

Файл workflow: `.github/workflows/ci.yml`

Pipeline запускается автоматически при:

- `push` в ветки `main` и `master`
- `pull_request`

В CI есть 4 отдельные job:

- `backend-build` — сборка backend
- `backend-test` — запуск backend-тестов с `gomock`, без подключения к БД
- `frontend-test` — запуск frontend-тестов
- `frontend-build` — сборка frontend

Job идут строго последовательно, без параллельного выполнения:

```text
backend-build -> backend-test -> frontend-test -> frontend-build
```

Если предыдущий этап падает, следующий не запускается.

## REST API

### Создать игру

`POST /api/games`

Пример тела запроса:

```json
{
  "title": "Cyberpunk 2077",
  "genre": "RPG",
  "platform": "PC",
  "releaseYear": 2020,
  "rating": 8,
  "status": "playing",
  "description": "Футуристическая RPG с открытым городом и вариативной прокачкой.",
  "imagePath": "/blob/cyberpunk-cover.jpg"
}
```

Пример:

```bash
curl -X POST http://localhost:3000/api/games \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Cyberpunk 2077",
    "genre": "RPG",
    "platform": "PC",
    "releaseYear": 2020,
    "rating": 8,
    "status": "playing",
    "description": "Футуристическая RPG с открытым городом и вариативной прокачкой.",
    "imagePath": "/blob/cyberpunk-cover.jpg"
  }'
```

### Получить список игр

`GET /api/games`

Пример:

```bash
curl http://localhost:3000/api/games
```

С фильтром:

```bash
curl "http://localhost:3000/api/games?status=completed&genre=RPG"
```

### Получить одну игру

`GET /api/games/:id`

```bash
curl http://localhost:3000/api/games/1
```

### Загрузить изображение

`POST /api/uploads/image`

Пример:

```bash
curl -X POST http://localhost:3000/api/uploads/image \
  -F "image=@cover.png"
```

Ответ:

```json
{
  "path": "/blob/1711111111111.png"
}
```

### Обновить игру

`PUT /api/games/:id`

```bash
curl -X PUT http://localhost:3000/api/games/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Cyberpunk 2077",
    "genre": "RPG",
    "platform": "PC",
    "releaseYear": 2020,
    "rating": 9,
    "status": "completed",
    "description": "Обновлённая версия игры после прохождения основных сюжетных линий.",
    "imagePath": "/blob/cyberpunk-cover.jpg"
  }'
```

### Удалить игру

`DELETE /api/games/:id`

```bash
curl -X DELETE http://localhost:3000/api/games/1
```

## Что проверяется тестами

Backend:

- создание игры
- получение списка игр
- обновление игры
- удаление игры
- обработка невалидных данных

Frontend:

- рендер списка игр
- отображение формы
- добавление игры
- удаление игры
- пользовательский сценарий открытия карточки игры
