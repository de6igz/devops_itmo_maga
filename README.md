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
├── backend
│   ├── app.go
│   ├── app_test.go
│   ├── go.mod
│   ├── internal
│   │   ├── game
│   │   ├── httpapi
│   │   ├── media
│   │   └── storage/postgres
│   ├── main.go
│   └── go.sum
├── frontend
│   ├── index.html
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

Для локального запуска удобно заранее создать базу:

```bash
createdb game_catalog
```

```bash
cd backend
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/game_catalog?sslmode=disable"
go mod tidy
go run .
```

Backend стартует на `http://localhost:3000`.

При первом запуске:

- автоматически создаётся таблица `games`
- автоматически добавляются несколько демонстрационных игр
- автоматически используется папка `backend/blob` для обложек игр

Доступные команды backend:

```bash
go run .
go build ./...
go test ./...
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
