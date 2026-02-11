# AI Agent Workflow Guide
## Организация разработки теннисной платформы с AI-агентами

---

## 1. Архитектура рабочего процесса

### Роли

```
Ты (Дирижёр / PM)
│
├── Claude Chat — Архитектор, PM, Code Reviewer
│   Задачи: планирование спринтов, обновление документации,
│   архитектурные решения, review кода, дизайн API
│
├── Claude Code — Backend Developer (терминал)
│   Задачи: Go код, миграции, sqlc, WebSocket, Docker, тесты
│   Инструкция: CLAUDE.md в корне репо
│
├── Cursor / Windsurf — Frontend Developer (IDE)
│   Задачи: Expo/RN мобайл, React web, компоненты, стили
│   Инструкция: .cursorrules в корне репо
│
└── Claude Code — DevOps (по необходимости)
    Задачи: Docker, CI/CD, деплой, мониторинг
```

### Главный принцип

**Документация = Source of Truth.** Агент НЕ принимает архитектурных решений. Агент реализует то, что описано в docs/. Если в документации нет — сначала обнови документ, потом давай задачу кодеру.

```
Идея / правка
  → Claude Chat обновляет документ
  → Ты коммитишь docs/ в репо
  → Агент-кодер читает docs/ и реализует
```

---

## 2. Структура репозитория

```
tennisapp/
├── CLAUDE.md                    ← Claude Code читает ПЕРВЫМ
├── .cursorrules                 ← Cursor читает ПЕРВЫМ
│
├── docs/                        ← Source of Truth
│   ├── 01-PRD.md
│   ├── 02-TECH-SPEC.md
│   ├── 03-DATABASE-SCHEMA.sql
│   ├── 04-API-SPECIFICATION.md
│   ├── 05-USER-STORIES.md
│   ├── 07-CODING-CONVENTIONS.md
│   ├── 08-ELO-ALGORITHM.md
│   ├── 09-INTEGRATIONS.md
│   ├── 10-DEPLOYMENT.md
│   ├── 11-SECURITY.md
│   ├── 12-TESTING-STRATEGY.md
│   ├── 13-DESIGN-SYSTEM.md
│   ├── mobile-app/              ← 10 спецификаций экранов
│   ├── web-platform/            ← 6 спецификаций web
│   └── sprints/                 ← Планы спринтов (создаются по ходу)
│       ├── sprint-01.md
│       └── sprint-02.md
│
├── apps/
│   ├── backend/                 ← Go
│   ├── mobile/                  ← Expo/RN
│   ├── web-admin/               ← React
│   └── web-superadmin/          ← React
│
├── packages/
│   ├── shared-types/
│   └── api-client/
│
├── scripts/
├── docker-compose.yml
└── Makefile
```

---

## 3. Подготовка к старту (чеклист)

```
□ Создать GitHub репозиторий
□ Скопировать все docs/ из project-v2/
□ Разместить CLAUDE.md в корне репо
□ Разместить .cursorrules в корне репо
□ Создать docker-compose.yml (Postgres + Redis + MinIO)
□ Создать Makefile
□ Создать .env.example
□ Инициализировать Go модуль (apps/backend/)
□ Инициализировать Expo проект (apps/mobile/)
□ Первый коммит + push
```

---

## 4. Спринтовый цикл

### Фаза 1: Планирование (30 мин) — Claude Chat

Ты говоришь:
> Мы начинаем Sprint N. Вот задачи из 05-USER-STORIES.md: [список].
> Разбей на конкретные шаги для агента-кодера.

Claude Chat выдаёт: порядок задач, промпты для агента, ожидаемые файлы. Ты сохраняешь в `docs/sprints/sprint-N.md`.

### Фаза 2: Backend — Claude Code

Начало каждой сессии:
```
Прочитай CLAUDE.md и docs/sprints/sprint-N.md.
Выполни задачу [ID]: [описание].
```

Правила:
- Одна задача за раз
- После каждой задачи: `go build ./...`
- Попроси написать тест
- Коммит после каждой задачи

### Фаза 3: Frontend — Cursor

Начало каждой сессии:
```
Контекст: Sprint N, экран [название].
Спецификация: docs/mobile-app/[файл].md
Backend API уже готов: [список эндпоинтов]
```

### Фаза 4: Проверка

1. Запусти `make dev` + `npx expo start`
2. Пройди user flow на устройстве
3. Запиши баги/правки

### Фаза 5: Правки

| Тип | Действие |
|-----|---------|
| Баг в коде | Опиши агенту: файл, функция, ошибка |
| UI косметика | Cursor: "Измени X в компоненте Y" |
| Новая фича | Claude Chat → обнови доки → агент кодит |
| Изменение API | Claude Chat → API Spec + Schema → backend → frontend |
| Архитектура | Claude Chat → PRD + Tech Spec → всё по цепочке |

### Фаза 6: Завершение спринта (15 мин) — Claude Chat

> Sprint N завершён. Сделано: [список]. Не успели: [список].
> Обнови sprint-N.md, подготовь план Sprint N+1.

---

## 5. Управление контекстом

### Что давать агенту

| Всегда | CLAUDE.md или .cursorrules |
|--------|---------------------------|
| По задаче | 1-2 конкретных документа (секция API, таблица БД, экран) |

### Что НЕ давать

- PRD целиком (слишком абстрактно)
- User Stories целиком (нужна 1 задача)
- Prototype JSX (только для визуального референса)

---

## 6. Sprint-by-Sprint план

### Sprint 1-2: Foundation + Auth (Недели 1-4)

```
Неделя 1: Backend foundation
  Claude Code → Go scaffold, Docker, Makefile, миграции, middleware

Неделя 2: Auth + Mobile shell
  Claude Code → Auth endpoints (OTP, JWT, profile, quiz, PIN)
  Cursor     → Expo init, design system, auth screens, API client
  Интеграция → Auth flow: phone → OTP → profile → quiz → home
```

### Sprint 3-4: Core Features (Недели 5-8)

```
  Claude Code → Users, Communities, Events CRUD
  Cursor     → Profile, Communities, Events экраны
  Порядок: backend endpoint → mobile экран → интеграция
```

### Sprint 5-6: Matches & Social (Недели 9-12)

```
  Claude Code → Match results, ELO (docs/08-ELO-ALGORITHM.md), WebSocket, Push
  Cursor     → Ввод результатов, чат, уведомления
```

### Sprint 7-8: Features (Недели 13-16)

```
  Claude Code → Posts, badges, friends, courts
  Cursor     → Home tab, посты, бейджи, карта кортов
```

### Sprint 9-10: Web Admin (Недели 17-20)

```
  Cursor     → Admin panel (React + Shadcn/UI)
  Cursor     → Superadmin panel
```

### Sprint 11-12: Polish & Launch (Недели 21-24)

```
  Все агенты → empty/error/loading states, seed data
  Claude Code → production deploy, CI/CD
  Ручное тестирование → баг-фиксы → beta launch
```

---

## 7. Промпт-шаблоны

### Backend задача
```
Прочитай CLAUDE.md.
Задача [ID]: [описание]
API: docs/04-API-SPECIFICATION.md, секция [X]
БД: docs/03-DATABASE-SCHEMA.sql, таблицы [X, Y]

Создай: handler → service → repository → dto
После: go build, напиши тест, запусти тесты
```

### Frontend задача
```
Контекст: Sprint N, экран [название]
Спецификация: docs/mobile-app/[файл].md
API endpoints: [список готовых]
Дизайн: docs/13-DESIGN-SYSTEM.md
Используй: TanStack Query, Zustand, компоненты из src/shared/ui/
```

### Code Review
```
Review файла [путь]. Проверь:
1. Соответствие docs/07-CODING-CONVENTIONS.md
2. Обработка ошибок
3. Security (SQL injection, validation)
4. N+1 запросы
5. Бизнес-логика по docs/01-PRD.md
```

### Обновление документации
```
Нужно добавить [описание].
Обнови: API Spec, Schema, User Stories.
Покажи только изменения (diff).
```

---

## 8. Git workflow

### Коммиты
```
feat(auth): implement OTP send
fix(chat): correct WebSocket reconnection
refactor(users): extract validation
test(rating): add ELO calculation tests
docs(api): update events specification
chore(ci): setup GitHub Actions
```

### Ветки
```
feature/auth-otp
feature/events-crud
fix/event-join-validation
```

---

## 9. Чеклист качества

### Backend
```
□ go build ./... — компилируется
□ go test ./... — тесты зелёные
□ Endpoints соответствуют API Spec
□ Ошибки в стандартном формате
□ Нет SELECT без LIMIT
□ Rate limiting по SECURITY.md
```

### Mobile
```
□ TypeScript strict — нет ошибок
□ Экран соответствует спецификации
□ Loading / Empty / Error states
□ Pull-to-refresh (списки)
□ iOS + Android работает
□ i18n ключи (ru, kk, en)
```

---

## 10. Обработка проблем

| Проблема | Решение |
|----------|---------|
| Агент создал файл не там | Укажи точный путь |
| Использовал не тот пакет | Укажи: "используй X, НЕ Y, см. Tech Spec" |
| Ломает существующий код | "Запусти go test ./... и покажи что упало" |
| Не знает как сделать | "Прочитай docs/09-INTEGRATIONS.md секция X" |
| Повторяет ошибку | Дай больше контекста из конкретного документа |
| Много файлов менять | Разбей на 5 мелких шагов |

---

## 11. Полезные команды

```bash
# Backend
make dev              # Dev server с hot-reload
make test             # Тесты
make migrate-up       # Миграции
make sqlc             # Генерация Go из SQL
make lint             # Линтер

# Mobile
npx expo start        # Dev server
npx expo start --ios  # iOS симулятор
eas build --platform ios  # Prod build

# Web
npm run dev           # Dev server (web-admin)
npm run build         # Prod build

# Infra
docker-compose up -d  # Postgres + Redis + MinIO
docker-compose logs -f # Логи
```
