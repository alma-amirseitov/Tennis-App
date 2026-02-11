# Структура файлов для разработки

## Монорепо: общая структура

```
tennisapp/
├── apps/
│   ├── mobile/                    # React Native (Expo)
│   ├── web-admin/                 # React (Community Admin Panel)
│   ├── web-superadmin/            # React (Superadmin Panel)
│   └── backend/                   # Go (API Server)
│
├── packages/
│   ├── shared-types/              # Общие TypeScript типы (mobile + web)
│   └── api-client/                # Общий API клиент
│
├── docs/                          # Документация (содержимое project-v2/)
├── scripts/                       # CI/CD, миграции, сиды
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## Mobile App (React Native / Expo)

```
apps/mobile/
├── app.json
├── package.json
├── tsconfig.json
├── babel.config.js
├── eas.json                        # EAS Build config
│
├── assets/
│   ├── images/                     # Иллюстрации, иконки
│   ├── fonts/                      # Кастомные шрифты
│   └── animations/                 # Lottie анимации
│
├── src/
│   ├── app/                        # Expo Router (file-based routing)
│   │   ├── _layout.tsx             # Root layout (auth guard, providers)
│   │   ├── index.tsx               # Splash / redirect
│   │   │
│   │   ├── (auth)/                 # Auth flow (без навигации)
│   │   │   ├── _layout.tsx
│   │   │   ├── onboarding.tsx
│   │   │   ├── phone.tsx           # Ввод телефона
│   │   │   ├── otp.tsx             # Ввод SMS-кода
│   │   │   ├── pin.tsx             # Установка / ввод PIN
│   │   │   ├── profile-setup.tsx   # Создание профиля
│   │   │   └── skill-quiz.tsx      # Опросник уровня
│   │   │
│   │   ├── (tabs)/                 # Main tab navigator
│   │   │   ├── _layout.tsx         # Tab bar config
│   │   │   ├── home.tsx            # Главная
│   │   │   ├── players.tsx         # Игроки
│   │   │   ├── events/             # Ивенты (nested)
│   │   │   │   ├── _layout.tsx     # Top tabs: Лента | Календарь | Мои
│   │   │   │   ├── feed.tsx
│   │   │   │   ├── calendar.tsx
│   │   │   │   └── mine.tsx
│   │   │   ├── communities.tsx     # Сообщества
│   │   │   └── profile.tsx         # Мой профиль
│   │   │
│   │   ├── event/
│   │   │   ├── [id].tsx            # Детальный экран ивента
│   │   │   ├── create.tsx          # Конструктор ивента (wizard)
│   │   │   └── result/[id].tsx     # Ввод результата
│   │   │
│   │   ├── community/
│   │   │   ├── [id].tsx            # Страница сообщества
│   │   │   ├── create.tsx          # Создание сообщества
│   │   │   └── manage/[id].tsx     # Админ-функции (в приложении)
│   │   │
│   │   ├── player/
│   │   │   └── [id].tsx            # Публичный профиль
│   │   │
│   │   ├── chat/
│   │   │   ├── index.tsx           # Список чатов
│   │   │   └── [id].tsx            # Экран чата
│   │   │
│   │   ├── notifications.tsx       # Экран уведомлений
│   │   ├── courts-map.tsx          # Карта кортов
│   │   ├── friends.tsx             # Список друзей
│   │   └── settings/
│   │       ├── index.tsx           # Главный экран настроек
│   │       ├── edit-profile.tsx
│   │       ├── notifications.tsx   # Настройки уведомлений
│   │       ├── privacy.tsx
│   │       ├── language.tsx
│   │       └── pin.tsx
│   │
│   ├── components/                 # UI компоненты
│   │   ├── ui/                     # Базовые (Button, Input, Card, Badge...)
│   │   │   ├── Button.tsx
│   │   │   ├── Input.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── Badge.tsx
│   │   │   ├── Avatar.tsx
│   │   │   ├── Chip.tsx
│   │   │   ├── Modal.tsx
│   │   │   ├── BottomSheet.tsx
│   │   │   ├── Skeleton.tsx
│   │   │   ├── EmptyState.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── layout/                 # Layout компоненты
│   │   │   ├── Header.tsx          # Глобальный хедер (💬 🔔)
│   │   │   ├── Screen.tsx          # Обёртка экрана (SafeArea + loading/error)
│   │   │   └── TabBar.tsx          # Кастомный tab bar
│   │   │
│   │   ├── home/                   # Компоненты Home таба
│   │   │   ├── RatingWidget.tsx
│   │   │   ├── QuickActions.tsx
│   │   │   ├── UpcomingGames.tsx
│   │   │   └── FeedPost.tsx
│   │   │
│   │   ├── events/                 # Компоненты Events
│   │   │   ├── EventCard.tsx
│   │   │   ├── EventFilters.tsx
│   │   │   ├── EventWizard/        # Многошаговый wizard
│   │   │   │   ├── StepType.tsx
│   │   │   │   ├── StepFormat.tsx
│   │   │   │   ├── StepMatchRules.tsx
│   │   │   │   ├── StepTournament.tsx
│   │   │   │   ├── StepLocation.tsx
│   │   │   │   ├── StepParticipants.tsx
│   │   │   │   ├── StepPricing.tsx
│   │   │   │   └── StepPreview.tsx
│   │   │   ├── TournamentBracket.tsx
│   │   │   ├── RoundRobinTable.tsx
│   │   │   ├── ScoreInput.tsx      # Ввод счёта
│   │   │   └── MatchResultCard.tsx
│   │   │
│   │   ├── players/
│   │   │   ├── PlayerCard.tsx
│   │   │   └── PlayerFilters.tsx
│   │   │
│   │   ├── communities/
│   │   │   ├── CommunityCard.tsx
│   │   │   ├── CommunityHeader.tsx
│   │   │   ├── MemberList.tsx
│   │   │   ├── Leaderboard.tsx
│   │   │   └── PhotoGallery.tsx
│   │   │
│   │   ├── profile/
│   │   │   ├── StatsWidget.tsx
│   │   │   ├── RatingChart.tsx     # График динамики рейтинга
│   │   │   ├── MatchHistory.tsx
│   │   │   ├── AchievementBadge.tsx
│   │   │   └── PostCard.tsx
│   │   │
│   │   └── chat/
│   │       ├── ChatListItem.tsx
│   │       ├── MessageBubble.tsx
│   │       └── ChatInput.tsx
│   │
│   ├── hooks/                      # Кастомные хуки
│   │   ├── useAuth.ts
│   │   ├── useWebSocket.ts
│   │   ├── useLocation.ts
│   │   ├── usePushNotifications.ts
│   │   ├── useRefreshToken.ts
│   │   └── useDebounce.ts
│   │
│   ├── services/                   # API и внешние сервисы
│   │   ├── api/                    # REST API клиент
│   │   │   ├── client.ts           # Axios/fetch instance + interceptors
│   │   │   ├── auth.ts
│   │   │   ├── events.ts
│   │   │   ├── communities.ts
│   │   │   ├── players.ts
│   │   │   ├── chat.ts
│   │   │   ├── notifications.ts
│   │   │   ├── courts.ts
│   │   │   └── ratings.ts
│   │   │
│   │   ├── websocket.ts            # WebSocket клиент для чата
│   │   ├── push.ts                 # Firebase FCM setup
│   │   └── storage.ts              # Secure storage (react-native-keychain)
│   │
│   ├── store/                      # State management (Zustand)
│   │   ├── authStore.ts
│   │   ├── chatStore.ts
│   │   ├── notificationStore.ts
│   │   └── uiStore.ts             # UI state (filters, modals)
│   │
│   ├── i18n/                       # Интернационализация
│   │   ├── index.ts                # i18next config
│   │   ├── ru.json
│   │   ├── kz.json
│   │   └── en.json
│   │
│   ├── utils/                      # Утилиты
│   │   ├── validation.ts           # Валидация форм
│   │   ├── formatters.ts           # Даты, числа, рейтинг
│   │   ├── constants.ts            # NTRP levels, districts, etc.
│   │   └── deeplink.ts             # Deep link handling
│   │
│   ├── types/                      # TypeScript типы
│   │   ├── user.ts
│   │   ├── event.ts
│   │   ├── community.ts
│   │   ├── chat.ts
│   │   ├── notification.ts
│   │   ├── rating.ts
│   │   └── navigation.ts
│   │
│   └── theme/                      # Дизайн-система
│       ├── colors.ts
│       ├── typography.ts
│       ├── spacing.ts
│       └── index.ts
```

---

## Web Admin Panel (React)

```
apps/web-admin/
├── package.json
├── tsconfig.json
├── vite.config.ts
│
├── public/
│   └── favicon.ico
│
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── routes.tsx                  # React Router config
│   │
│   ├── pages/
│   │   ├── auth/
│   │   │   └── Login.tsx
│   │   │
│   │   ├── dashboard/
│   │   │   └── Dashboard.tsx
│   │   │
│   │   ├── members/
│   │   │   ├── MemberList.tsx
│   │   │   ├── MemberDetail.tsx
│   │   │   └── Applications.tsx    # Заявки на вступление
│   │   │
│   │   ├── events/
│   │   │   ├── EventList.tsx
│   │   │   ├── EventCreate.tsx
│   │   │   ├── EventDetail.tsx
│   │   │   └── tournament/
│   │   │       ├── BracketView.tsx
│   │   │       ├── RoundRobin.tsx
│   │   │       ├── DrawGenerator.tsx
│   │   │       └── ResultEntry.tsx
│   │   │
│   │   ├── content/
│   │   │   ├── PostList.tsx
│   │   │   ├── PostEditor.tsx
│   │   │   └── PhotoGallery.tsx
│   │   │
│   │   ├── rating/
│   │   │   ├── Leaderboard.tsx
│   │   │   └── RatingSettings.tsx
│   │   │
│   │   ├── finances/
│   │   │   ├── Overview.tsx
│   │   │   ├── Transactions.tsx
│   │   │   └── Subscriptions.tsx
│   │   │
│   │   ├── courts/                 # Только для клубов
│   │   │   └── CourtSchedule.tsx
│   │   │
│   │   └── settings/
│   │       ├── General.tsx
│   │       ├── Access.tsx
│   │       └── DangerZone.tsx
│   │
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Sidebar.tsx
│   │   │   ├── Header.tsx
│   │   │   ├── PageContainer.tsx
│   │   │   └── CommunitySwitch.tsx # Переключение между сообществами
│   │   │
│   │   ├── ui/                     # Shadcn/UI wrappers
│   │   │   └── ...
│   │   │
│   │   ├── charts/                 # Графики (recharts)
│   │   │   ├── GrowthChart.tsx
│   │   │   ├── ActivityChart.tsx
│   │   │   └── LevelDistribution.tsx
│   │   │
│   │   └── data-table/             # Переиспользуемая таблица
│   │       ├── DataTable.tsx
│   │       ├── ColumnHeader.tsx
│   │       └── Pagination.tsx
│   │
│   ├── hooks/
│   ├── services/
│   ├── store/
│   ├── types/
│   ├── utils/
│   └── i18n/
```

---

## Web Superadmin Panel

```
apps/web-superadmin/
├── src/
│   ├── pages/
│   │   ├── dashboard/              # Platform-wide metrics
│   │   ├── users/                  # All users management
│   │   ├── communities/            # Communities + verification queue
│   │   ├── events/                 # All events moderation
│   │   ├── courts/                 # CRUD for courts
│   │   ├── notifications/          # Mass push notifications
│   │   ├── finances/               # Platform-wide finances
│   │   ├── analytics/              # Cohorts, funnels, retention
│   │   └── system/                 # Config, feature flags, logs
│   └── ...                         # Same structure as web-admin
```

---

## Backend (Go)

```
apps/backend/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── config/
│   │   └── config.go               # Env vars, settings
│   │
│   ├── handler/                    # HTTP handlers (controllers)
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── event.go
│   │   ├── community.go
│   │   ├── chat.go
│   │   ├── notification.go
│   │   ├── rating.go
│   │   ├── court.go
│   │   ├── post.go
│   │   ├── admin.go                # Community admin endpoints
│   │   ├── superadmin.go           # Superadmin endpoints
│   │   └── middleware/
│   │       ├── auth.go             # JWT validation
│   │       ├── rbac.go             # Role-based access
│   │       ├── ratelimit.go
│   │       └── cors.go
│   │
│   ├── service/                    # Business logic
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── event_service.go
│   │   ├── community_service.go
│   │   ├── chat_service.go
│   │   ├── notification_service.go
│   │   ├── rating_service.go       # ELO calculation
│   │   ├── court_service.go
│   │   ├── tournament_service.go   # Bracket generation, draws
│   │   ├── sms_service.go          # SMS OTP sending
│   │   └── push_service.go         # Firebase FCM
│   │
│   ├── repository/                 # Database layer
│   │   ├── user_repo.go
│   │   ├── event_repo.go
│   │   ├── community_repo.go
│   │   ├── chat_repo.go
│   │   ├── notification_repo.go
│   │   ├── rating_repo.go
│   │   ├── court_repo.go
│   │   └── post_repo.go
│   │
│   ├── model/                      # Domain models
│   │   ├── user.go
│   │   ├── event.go
│   │   ├── community.go
│   │   ├── match.go
│   │   ├── rating.go
│   │   ├── chat.go
│   │   ├── notification.go
│   │   ├── court.go
│   │   └── post.go
│   │
│   ├── ws/                         # WebSocket
│   │   ├── hub.go                  # Connection manager
│   │   ├── client.go               # Client connection
│   │   └── message.go              # Message types
│   │
│   └── pkg/                        # Shared utilities
│       ├── jwt/
│       ├── sms/                    # SMS provider adapter
│       ├── firebase/               # FCM client
│       ├── storage/                # File upload (S3/MinIO)
│       ├── elo/                    # Rating algorithm
│       └── validator/
│
├── migrations/                     # SQL migrations
│   ├── 001_users.up.sql
│   ├── 001_users.down.sql
│   ├── 002_communities.up.sql
│   ├── 003_events.up.sql
│   ├── 004_matches.up.sql
│   ├── 005_ratings.up.sql
│   ├── 006_chats.up.sql
│   ├── 007_notifications.up.sql
│   ├── 008_courts.up.sql
│   ├── 009_posts.up.sql
│   └── 010_achievements.up.sql
│
├── seeds/                          # Test data
│   ├── users.sql
│   ├── communities.sql
│   ├── courts_astana.sql           # Реальные корты Астаны
│   └── achievements.sql
│
├── Dockerfile
├── go.mod
└── go.sum
```

---

## Docker Compose (Development)

```yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: tennisapp
      POSTGRES_USER: tennis
      POSTGRES_PASSWORD: dev_password
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build: ./apps/backend
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    environment:
      DATABASE_URL: postgres://tennis:dev_password@postgres:5432/tennisapp
      REDIS_URL: redis://redis:6379

  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin

volumes:
  pgdata:
```
