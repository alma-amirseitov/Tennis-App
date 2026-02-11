# API Specification v2.0
## Tennis Platform API

**Base URL:** `https://api.tennisapp.kz/v1`  
**Format:** REST JSON  
**Auth:** Bearer JWT  
**Date:** 2026-02-10

---

## Conventions

**Auth Header:** `Authorization: Bearer <access_token>`

**Success Response:**
```json
{ "data": { ... } }
```

**List Response:**
```json
{
  "data": [ ... ],
  "pagination": { "page": 1, "per_page": 20, "total": 142, "total_pages": 8 }
}
```

**Error Response:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Описание ошибки",
    "details": [{ "field": "phone", "message": "Неверный формат" }]
  }
}
```

**Error Codes:**
| HTTP | Code | Описание |
|------|------|----------|
| 400 | VALIDATION_ERROR | Невалидные данные |
| 401 | UNAUTHORIZED | Нет токена / истёк |
| 403 | FORBIDDEN | Нет доступа |
| 404 | NOT_FOUND | Ресурс не найден |
| 409 | CONFLICT | Конфликт (уже существует) |
| 429 | RATE_LIMITED | Превышен лимит |
| 500 | INTERNAL_ERROR | Ошибка сервера |

---

## 1. AUTH (5 endpoints)

### POST /auth/otp/send
Отправить SMS с OTP-кодом.

**Rate limit:** 3/час, 10/день на номер

**Request:**
```json
{ "phone": "+77071234567" }
```

**Response 200:**
```json
{
  "data": {
    "session_id": "uuid",
    "expires_in": 300,
    "retry_after": 60
  }
}
```

---

### POST /auth/otp/verify
Проверить OTP-код. Возвращает токены если пользователь существует, или temp_token для нового.

**Rate limit:** 5 попыток на session

**Request:**
```json
{
  "session_id": "uuid",
  "code": "1234"
}
```

**Response 200 (existing user):**
```json
{
  "data": {
    "is_new": false,
    "access_token": "jwt...",
    "refresh_token": "rt...",
    "user": { "id": "uuid", "first_name": "Иван", "is_profile_complete": true }
  }
}
```

**Response 200 (new user):**
```json
{
  "data": {
    "is_new": true,
    "temp_token": "jwt...",
    "user_id": "uuid"
  }
}
```

---

### POST /auth/profile/setup 🔒 temp_token
Заполнить профиль нового пользователя. Обязательный шаг после регистрации.

**Request:**
```json
{
  "first_name": "Иван",
  "last_name": "Петров",
  "gender": "male",
  "birth_year": 1990,
  "city": "Астана",
  "district": "Есильский",
  "language": "ru"
}
```

**Response 200:**
```json
{
  "data": {
    "access_token": "jwt...",
    "refresh_token": "rt...",
    "user": { ... }
  }
}
```

---

### POST /auth/refresh
Обновить access token через refresh token.

**Request:**
```json
{ "refresh_token": "rt..." }
```

**Response 200:**
```json
{
  "data": {
    "access_token": "jwt...",
    "refresh_token": "rt_new..."
  }
}
```

---

### POST /auth/pin/set 🔒
Установить или обновить PIN-код.

**Request:**
```json
{ "pin": "1234" }
```

**Response 200:**
```json
{ "data": { "message": "PIN установлен" } }
```

---

### POST /auth/pin/verify
Войти по PIN (без SMS).

**Rate limit:** 3 попытки, потом блок на 15 мин

**Request:**
```json
{
  "user_id": "uuid",
  "pin": "1234"
}
```

**Response 200:** (same as otp/verify for existing user)

---

## 2. QUIZ (2 endpoints)

### GET /quiz/questions 🔒
Получить вопросы skill quiz.

**Response 200:**
```json
{
  "data": {
    "questions": [
      {
        "id": 1,
        "text_ru": "Как давно вы играете в теннис?",
        "text_kz": "...",
        "text_en": "...",
        "options": [
          { "id": "a", "text_ru": "Менее 1 года", "weight": 1 },
          { "id": "b", "text_ru": "1-3 года", "weight": 2 },
          { "id": "c", "text_ru": "3-5 лет", "weight": 3 },
          { "id": "d", "text_ru": "Более 5 лет", "weight": 4 }
        ]
      }
    ]
  }
}
```

---

### POST /quiz/submit 🔒
Отправить ответы quiz. Возвращает определённый уровень.

**Request:**
```json
{
  "answers": [
    { "question_id": 1, "option_id": "b" },
    { "question_id": 2, "option_id": "c" },
    { "question_id": 3, "option_id": "b" }
  ]
}
```

**Response 200:**
```json
{
  "data": {
    "ntrp_level": 3.0,
    "level_label": "Любитель",
    "initial_rating": 1200.00
  }
}
```

---

## 3. USERS (8 endpoints)

### GET /users/me 🔒
Получить свой профиль.

**Response 200:**
```json
{
  "data": {
    "id": "uuid",
    "phone": "+7707***4567",
    "first_name": "Иван",
    "last_name": "Петров",
    "gender": "male",
    "birth_year": 1990,
    "city": "Астана",
    "district": "Есильский",
    "avatar_url": "https://...",
    "bio": "Играю по выходным",
    "ntrp_level": 3.0,
    "level_label": "Любитель",
    "global_rating": 1250.50,
    "global_games_count": 42,
    "language": "ru",
    "is_profile_complete": true,
    "quiz_completed": true,
    "pin_set": true,
    "notification_settings": { ... },
    "created_at": "2026-01-15T10:00:00Z"
  }
}
```

---

### PATCH /users/me 🔒
Обновить свой профиль. Partial update.

**Request:**
```json
{
  "first_name": "Иван",
  "bio": "Играю 3 раза в неделю",
  "district": "Сарыаркинский",
  "language": "kz"
}
```

---

### POST /users/me/avatar 🔒
Загрузить аватар. Multipart form-data.

**Request:** `multipart/form-data`, field `avatar` (max 5MB, jpg/png/webp)

**Response 200:**
```json
{ "data": { "avatar_url": "https://storage.../avatar.jpg" } }
```

---

### PATCH /users/me/notifications 🔒
Обновить настройки уведомлений.

**Request:**
```json
{
  "event_response": true,
  "game_reminder_24h": false,
  "quiet_hours_start": "22:00",
  "quiet_hours_end": "08:00"
}
```

---

### PATCH /users/me/privacy 🔒
Обновить настройки приватности.

**Request:**
```json
{
  "profile_visibility": "communities",
  "allow_messages_from": "friends",
  "show_stats": false
}
```

---

### GET /users/:id 🔒
Публичный профиль пользователя.

**Response 200:**
```json
{
  "data": {
    "id": "uuid",
    "first_name": "Алексей",
    "last_name": "С.",
    "avatar_url": "...",
    "ntrp_level": 3.5,
    "level_label": "Продвинутый любитель",
    "global_rating": 1400.00,
    "stats": {
      "total_games": 85,
      "total_wins": 52,
      "win_rate": 61.2,
      "current_streak": 3
    },
    "badges": [
      { "id": "ten_wins", "icon": "🎖", "name": "Десятка", "earned_at": "..." }
    ],
    "communities": [
      { "id": "uuid", "name": "NTC Astana", "role": "member" }
    ],
    "is_friend": false,
    "mutual_communities": 2
  }
}
```

---

### GET /users/search 🔒
Поиск игроков с фильтрами.

**Query params:**
| Param | Type | Description |
|-------|------|-------------|
| q | string | Поиск по имени |
| min_level | float | Мин. NTRP |
| max_level | float | Макс. NTRP |
| gender | string | male/female |
| district | string | Район |
| community_id | uuid | Участники сообщества |
| sort | string | rating / activity / name / games |
| page | int | Страница |
| per_page | int | Кол-во (max 50) |

**Response 200:** List of user profiles with pagination.

---

### POST /users/me/push-token 🔒
Обновить FCM push token.

**Request:**
```json
{ "token": "fcm_token_string", "device": "iPhone 15" }
```

---

## 4. FRIENDS (4 endpoints)

### GET /friends 🔒
Список друзей.

**Response 200:**
```json
{
  "data": [
    {
      "id": "uuid",
      "user": { "id": "uuid", "first_name": "Алексей", "avatar_url": "...", "ntrp_level": 3.5 },
      "added_at": "2026-01-20T15:00:00Z"
    }
  ]
}
```

---

### POST /friends/:user_id 🔒
Добавить в друзья (одностороннее).

**Response 201:**
```json
{ "data": { "message": "Добавлен в друзья" } }
```

---

### DELETE /friends/:user_id 🔒
Удалить из друзей.

---

### GET /friends/check/:user_id 🔒
Проверить является ли пользователь другом.

**Response 200:**
```json
{ "data": { "is_friend": true } }
```

---

## 5. EVENTS (14 endpoints)

### GET /events 🔒
Лента ивентов с фильтрами.

**Query params:**
| Param | Type | Description |
|-------|------|-------------|
| type | string | find_partner / organized_game / tournament / training |
| status | string | published / registration_open / in_progress / completed |
| composition | string | singles / doubles / mixed |
| community_id | uuid | Ивенты сообщества |
| min_level | float | Мин. уровень |
| max_level | float | Макс. уровень |
| date_from | date | С даты |
| date_to | date | По дату |
| district | string | Район |
| sort | string | date_asc / date_desc / spots_left |
| page, per_page | int | Пагинация |

---

### GET /events/calendar 🔒
Ивенты для календаря (сгруппированы по дням).

**Query params:** `month` (2026-03), `community_id` (optional)

**Response 200:**
```json
{
  "data": {
    "2026-03-15": [
      { "id": "uuid", "title": "...", "start_time": "...", "event_type": "...", "status": "..." }
    ],
    "2026-03-16": [ ... ]
  }
}
```

---

### GET /events/my 🔒
Мои ивенты (созданные + записанные).

**Query params:** `tab` (created / joined / past), `page`, `per_page`

---

### GET /events/:id 🔒
Детали ивента.

**Response 200:**
```json
{
  "data": {
    "id": "uuid",
    "title": "Вечерняя парная игра",
    "description": "...",
    "event_type": "organized_game",
    "status": "registration_open",
    "community": { "id": "uuid", "name": "NTC Astana", "logo_url": "..." },
    "player_composition": "doubles",
    "match_format": "best_of",
    "match_format_details": { "sets": 3, "games_per_set": 6, "tiebreak": true },
    "court": { "id": "uuid", "name": "NTC Astana", "address": "..." },
    "start_time": "2026-03-15T18:00:00+06:00",
    "end_time": "2026-03-15T20:00:00+06:00",
    "max_participants": 8,
    "current_participants": 5,
    "min_level": 2.5,
    "max_level": 4.0,
    "registration_deadline": "2026-03-15T16:00:00+06:00",
    "is_paid": false,
    "created_by": { "id": "uuid", "first_name": "Алексей" },
    "participants": [
      { "id": "uuid", "first_name": "Иван", "avatar_url": "...", "ntrp_level": 3.0, "status": "registered" }
    ],
    "my_status": "registered",
    "can_join": false,
    "can_edit": true,
    "created_at": "..."
  }
}
```

---

### POST /events 🔒
Создать ивент (конструктор, 1 запрос).

**Request:**
```json
{
  "title": "Ищу партнёра на вечер",
  "description": "Уровень 3.0-4.0, 2 сета",
  "event_type": "find_partner",
  "community_id": null,
  "player_composition": "singles",
  "match_format": "best_of",
  "match_format_details": { "sets": 2, "games_per_set": 6, "tiebreak": true },
  "court_id": "uuid",
  "location_name": null,
  "location_address": null,
  "start_time": "2026-03-15T18:00:00+06:00",
  "end_time": "2026-03-15T20:00:00+06:00",
  "max_participants": 2,
  "min_level": 3.0,
  "max_level": 4.0,
  "gender_restriction": null,
  "registration_deadline": "2026-03-15T16:00:00+06:00",
  "status": "published"
}
```

---

### PATCH /events/:id 🔒
Обновить ивент. Только автор или админ сообщества.

---

### DELETE /events/:id 🔒
Удалить ивент (только draft/cancelled). Или отменить (published → cancelled).

---

### POST /events/:id/join 🔒
Записаться на ивент.

**Request (optional for doubles):**
```json
{ "partner_id": "uuid" }
```

**Response 201:**
```json
{ "data": { "participant_id": "uuid", "status": "registered" } }
```

**Error 400:** `ALREADY_JOINED`, `EVENT_FULL`, `LEVEL_MISMATCH`, `REGISTRATION_CLOSED`

---

### DELETE /events/:id/join 🔒
Отписаться от ивента.

---

### PATCH /events/:id/status 🔒
Изменить статус ивента (для автора / админа).

**Request:**
```json
{ "status": "in_progress" }
```

Allowed transitions:
- draft → published
- published → registration_open → registration_closed → in_progress → completed
- any → cancelled

---

### GET /events/:id/participants 🔒
Список участников ивента.

---

### POST /events/:id/matches 🔒
Создать матч внутри ивента (для организатора).

**Request:**
```json
{
  "player1_id": "uuid",
  "player2_id": "uuid",
  "composition": "singles",
  "round_name": "Полуфинал",
  "round_number": 2,
  "court_number": 1,
  "scheduled_time": "2026-03-15T18:00:00+06:00"
}
```

---

### GET /events/:id/matches 🔒
Список матчей ивента.

---

### GET /events/:id/bracket 🔒
Турнирная сетка (Phase 2, но endpoint зарезервирован).

---

## 6. MATCHES (5 endpoints)

### GET /matches/my 🔒
Моя история матчей.

**Query params:** `community_id`, `opponent_id`, `result` (win/loss/all), `page`, `per_page`

---

### POST /matches/:id/result 🔒
Внести результат матча.

**Request:**
```json
{
  "winner_id": "uuid",
  "score": {
    "sets": [
      { "p1": 6, "p2": 4 },
      { "p1": 3, "p2": 6 },
      { "p1": 7, "p2": 5, "tiebreak": { "p1": 7, "p2": 3 } }
    ]
  }
}
```

**Response 200:**
```json
{
  "data": {
    "match_id": "uuid",
    "result_status": "pending",
    "message": "Ожидает подтверждения от соперника"
  }
}
```

---

### POST /matches/:id/confirm 🔒
Подтвердить результат матча (второй игрок).

**Request:**
```json
{ "action": "confirm" }
```
или
```json
{ "action": "dispute", "reason": "Счёт был 6-4, 6-3, а не 6-4, 3-6, 7-5" }
```

**Response 200 (confirmed):**
```json
{
  "data": {
    "result_status": "confirmed",
    "rating_changes": {
      "player1": { "before": 1200.0, "after": 1218.5, "change": +18.5 },
      "player2": { "before": 1350.0, "after": 1331.5, "change": -18.5 }
    }
  }
}
```

---

### POST /matches/:id/admin-confirm 🔒 admin
Админ подтверждает спорный результат.

**Request:**
```json
{
  "winner_id": "uuid",
  "score": { ... },
  "note": "Подтверждено на основании записи камеры"
}
```

---

### GET /matches/:id 🔒
Детали матча.

---

## 7. COMMUNITIES (12 endpoints)

### GET /communities 🔒
Список сообществ.

**Query params:**
| Param | Type | Description |
|-------|------|-------------|
| type | string | club / league / organizer / group |
| access | string | open / closed |
| verified | bool | Только верифицированные |
| q | string | Поиск по названию |
| district | string | Район |
| sort | string | members / activity / name / created |
| page, per_page | int | Пагинация |

---

### GET /communities/my 🔒
Мои сообщества (где я участник).

---

### GET /communities/:id 🔒
Детали сообщества.

**Response 200:**
```json
{
  "data": {
    "id": "uuid",
    "name": "NTC Astana",
    "slug": "ntc-astana",
    "description": "...",
    "community_type": "club",
    "access_level": "open",
    "verification_status": "verified",
    "logo_url": "...",
    "banner_url": "...",
    "contact_phone": "+7...",
    "social_links": { "instagram": "...", "telegram": "..." },
    "address": "Кабанбай батыра, 42",
    "member_count": 245,
    "event_count": 18,
    "my_role": "member",
    "my_status": "active",
    "created_at": "..."
  }
}
```

---

### POST /communities 🔒
Создать сообщество.

**Request:**
```json
{
  "name": "Weekend Tennis Group",
  "description": "Играем по выходным в Есильском районе",
  "community_type": "group",
  "access_level": "open",
  "district": "Есильский"
}
```

Для club/league/organizer: verification_status автоматически ставится `pending`.

---

### PATCH /communities/:id 🔒 owner/admin
Обновить настройки сообщества.

---

### POST /communities/:id/join 🔒
Вступить / подать заявку.

**Request (для закрытых):**
```json
{ "message": "Хочу вступить, играю 3.0 NTRP" }
```

**Response 200 (open):**
```json
{ "data": { "status": "active", "message": "Вы вступили в сообщество" } }
```

**Response 200 (closed):**
```json
{ "data": { "status": "pending", "message": "Заявка отправлена" } }
```

---

### DELETE /communities/:id/join 🔒
Выйти из сообщества.

---

### GET /communities/:id/members 🔒
Список участников.

**Query params:** `role`, `status`, `q` (search), `sort`, `page`, `per_page`

---

### PATCH /communities/:id/members/:user_id 🔒 owner/admin
Изменить роль или статус участника.

**Request:**
```json
{ "role": "moderator" }
```
или
```json
{ "status": "banned", "reason": "Нарушение правил" }
```

---

### POST /communities/:id/members/review 🔒 admin/moderator
Одобрить или отклонить заявку.

**Request:**
```json
{
  "user_id": "uuid",
  "action": "approve"
}
```

---

### GET /communities/:id/leaderboard 🔒
Рейтинг участников сообщества.

**Query params:** `page`, `per_page`

**Response 200:**
```json
{
  "data": [
    {
      "rank": 1,
      "user": { "id": "uuid", "first_name": "Алексей", "avatar_url": "..." },
      "rating": 1450.50,
      "games": 32,
      "wins": 24,
      "losses": 8,
      "win_rate": 75.0
    }
  ]
}
```

---

### GET /communities/:id/feed 🔒
Лента постов сообщества.

**Query params:** `page`, `per_page`

---

## 8. POSTS (5 endpoints)

### GET /feed 🔒
Глобальная лента (посты сообществ + друзей + результаты матчей).

**Query params:** `tab` (news / feed), `page`, `per_page`

---

### POST /posts 🔒
Создать пост.

**Request:** `multipart/form-data`
| Field | Type | Description |
|-------|------|-------------|
| content | string | Текст поста |
| community_id | uuid | Если от имени сообщества (нужна роль admin+) |
| photos[] | file | До 5 фото (max 5MB каждое) |

---

### DELETE /posts/:id 🔒
Удалить пост (автор или админ сообщества).

---

### POST /posts/:id/like 🔒
Лайкнуть пост.

---

### DELETE /posts/:id/like 🔒
Убрать лайк.

---

## 9. CHAT (7 endpoints)

### GET /chats 🔒
Список чатов пользователя.

**Response 200:**
```json
{
  "data": [
    {
      "id": "uuid",
      "chat_type": "personal",
      "other_user": { "id": "uuid", "first_name": "Алексей", "avatar_url": "..." },
      "last_message": { "content": "Привет, играем завтра?", "sender_id": "uuid", "created_at": "..." },
      "unread_count": 2,
      "is_muted": false
    },
    {
      "id": "uuid",
      "chat_type": "community",
      "community": { "id": "uuid", "name": "NTC Astana", "logo_url": "..." },
      "last_message": { ... },
      "unread_count": 15,
      "is_muted": true
    }
  ]
}
```

---

### POST /chats/personal 🔒
Создать или получить существующий личный чат.

**Request:**
```json
{ "user_id": "uuid" }
```

**Response 200:**
```json
{ "data": { "chat_id": "uuid", "is_new": false } }
```

---

### GET /chats/:id/messages 🔒
Сообщения чата.

**Query params:** `before` (cursor, message_id), `limit` (default 50)

**Response 200:**
```json
{
  "data": [
    {
      "id": "uuid",
      "sender": { "id": "uuid", "first_name": "Иван", "avatar_url": "..." },
      "content": "Привет!",
      "reply_to": null,
      "created_at": "2026-03-15T18:30:00Z"
    }
  ],
  "has_more": true
}
```

---

### POST /chats/:id/messages 🔒
Отправить сообщение (HTTP fallback; основной путь — WebSocket).

**Request:**
```json
{
  "content": "Привет! Играем завтра?",
  "reply_to_id": null
}
```

---

### POST /chats/:id/read 🔒
Отметить чат как прочитанный.

**Request:**
```json
{ "last_read_at": "2026-03-15T18:35:00Z" }
```

---

### PATCH /chats/:id/mute 🔒
Замутить/размутить чат.

**Request:**
```json
{ "is_muted": true }
```

---

### GET /chats/unread-count 🔒
Общее количество непрочитанных сообщений (для badge).

**Response 200:**
```json
{ "data": { "total_unread": 17 } }
```

---

## 10. NOTIFICATIONS (4 endpoints)

### GET /notifications 🔒
Список уведомлений.

**Query params:** `page`, `per_page`

**Response 200:**
```json
{
  "data": [
    {
      "id": "uuid",
      "type": "result_confirm",
      "title": "Подтвердите результат",
      "body": "Алексей внёс результат: 6-4, 6-3",
      "data": { "match_id": "uuid", "deeplink": "tennis://matches/uuid" },
      "is_read": false,
      "created_at": "2026-03-15T20:00:00Z"
    }
  ]
}
```

---

### POST /notifications/read 🔒
Отметить уведомления как прочитанные.

**Request:**
```json
{ "ids": ["uuid1", "uuid2"] }
```
или
```json
{ "read_all": true }
```

---

### GET /notifications/unread-count 🔒
Количество непрочитанных (для badge 🔔).

**Response 200:**
```json
{ "data": { "count": 5 } }
```

---

### DELETE /notifications/:id 🔒
Удалить уведомление.

---

## 11. RATING (4 endpoints)

### GET /rating/global 🔒
Глобальный рейтинг.

**Query params:** `page`, `per_page`

**Response 200:**
```json
{
  "data": [
    {
      "rank": 1,
      "user": { "id": "uuid", "first_name": "Марат", "avatar_url": "...", "ntrp_level": 4.5 },
      "rating": 1650.00,
      "games": 120,
      "win_rate": 72.5
    }
  ]
}
```

---

### GET /rating/history 🔒
Моя история рейтинга (для графика).

**Query params:** `community_id` (null = global), `period` (1m / 3m / 6m / 1y / all)

**Response 200:**
```json
{
  "data": [
    { "date": "2026-01-15", "rating": 1000.00 },
    { "date": "2026-01-22", "rating": 1025.00 },
    { "date": "2026-02-01", "rating": 1080.50 }
  ]
}
```

---

### GET /rating/me 🔒
Моя позиция в рейтингах.

**Response 200:**
```json
{
  "data": {
    "global": { "rating": 1250.50, "rank": 45, "total_players": 500 },
    "communities": [
      { "community_id": "uuid", "community_name": "NTC Astana", "rating": 1300.00, "rank": 12 }
    ]
  }
}
```

---

### GET /rating/badges 🔒
Мои достижения и прогресс.

**Response 200:**
```json
{
  "data": {
    "earned": [
      { "id": "first_win", "icon": "🏅", "name": "Первая победа", "earned_at": "..." }
    ],
    "in_progress": [
      { "id": "ten_wins", "icon": "🎖", "name": "Десятка", "current": 7, "target": 10, "progress": 70 }
    ]
  }
}
```

---

## 12. COURTS (4 endpoints)

### GET /courts 🔒
Список кортов (для карты и выбора в конструкторе).

**Query params:** `district`, `surface`, `indoor` (bool), `q` (search)

**Response 200:**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "NTC Astana",
      "address": "Кабанбай батыра, 42",
      "district": "Есильский",
      "latitude": 51.1282,
      "longitude": 71.4307,
      "total_courts": 6,
      "indoor_courts": 4,
      "outdoor_courts": 2,
      "surface": "hard",
      "price_per_hour": 3000,
      "currency": "KZT",
      "phone": "+77172123456",
      "photos": ["url1", "url2"]
    }
  ]
}
```

---

### GET /courts/:id 🔒
Детали корта.

---

### GET /courts/map 🔒
Корты для карты (облегчённый формат).

**Response 200:**
```json
{
  "data": [
    { "id": "uuid", "name": "NTC Astana", "lat": 51.1282, "lng": 71.4307, "courts": 6, "surface": "hard" }
  ]
}
```

---

### POST /courts 🔒 superadmin
Создать корт (только суперадмин).

---

## 13. ADMIN — Web Panel (10 endpoints)

Все endpoint'ы требуют роль `owner`, `admin`, или `moderator` в сообществе.

### GET /admin/communities/:id/dashboard 🔒 admin
Дашборд сообщества.

**Response 200:**
```json
{
  "data": {
    "member_count": 245,
    "member_growth_30d": 12,
    "active_events": 3,
    "matches_this_month": 28,
    "top_players": [ ... ],
    "recent_activity": [ ... ]
  }
}
```

---

### GET /admin/communities/:id/members 🔒 admin
Список участников (таблица с фильтрами).

**Query params:** `status`, `role`, `q`, `sort`, `page`, `per_page`

---

### PATCH /admin/communities/:id/members/:user_id 🔒 admin
Изменить роль/статус (бан, повышение, понижение).

---

### POST /admin/communities/:id/members/review 🔒 admin
Bulk review заявок.

**Request:**
```json
{
  "actions": [
    { "user_id": "uuid1", "action": "approve" },
    { "user_id": "uuid2", "action": "reject", "reason": "Не теннисист" }
  ]
}
```

---

### POST /admin/communities/:id/events 🔒 admin
Создать ивент от имени сообщества.

---

### POST /admin/communities/:id/matches/:match_id/result 🔒 coach_referee+
Ввести результат матча (админ/судья — без подтверждения 2-й стороной).

---

### POST /admin/communities/:id/posts 🔒 admin
Создать пост от имени сообщества.

---

### GET /admin/communities/:id/export 🔒 admin
Экспорт данных (CSV).

**Query params:** `type` (members / matches / ratings)

**Response:** CSV file download.

---

### PATCH /admin/communities/:id/settings 🔒 owner
Обновить настройки сообщества (рейтинг, доступ, опасная зона).

---

### PATCH /admin/communities/:id/rating-settings 🔒 owner/admin
Настройки рейтинговой системы сообщества.

**Request:**
```json
{
  "initial_rating": 1000,
  "k_factor": 32,
  "min_games_for_leaderboard": 3
}
```

---

## 14. SUPERADMIN (7 endpoints)

Все endpoint'ы требуют `platform_role = superadmin`.

### GET /superadmin/dashboard 🔒 superadmin
Статистика платформы.

**Response 200:**
```json
{
  "data": {
    "total_users": 2400,
    "active_users_30d": 450,
    "total_communities": 15,
    "total_matches": 1200,
    "registrations_today": 5,
    "matches_today": 12
  }
}
```

---

### GET /superadmin/verifications 🔒 superadmin
Очередь верификации сообществ.

---

### POST /superadmin/verifications/:community_id 🔒 superadmin
Одобрить / отклонить верификацию.

**Request:**
```json
{ "action": "verify", "note": "Документы проверены" }
```

---

### POST /superadmin/users/:id/ban 🔒 superadmin
Забанить пользователя.

**Request:**
```json
{ "reason": "Спам", "duration_days": null }
```

---

### POST /superadmin/users/:id/unban 🔒 superadmin
Разбанить пользователя.

---

### CRUD /superadmin/courts 🔒 superadmin
CRUD кортов (POST/PATCH/DELETE).

---

### GET /superadmin/users 🔒 superadmin
Список всех пользователей (для модерации).

**Query params:** `status`, `q`, `sort`, `page`, `per_page`

---

## 15. WEBSOCKET

### Connection
```
wss://api.tennisapp.kz/ws?token={jwt_access_token}
```

### Client → Server messages

**Send message:**
```json
{
  "type": "message",
  "chat_id": "uuid",
  "content": "Привет!",
  "reply_to": null,
  "client_id": "temp-uuid"
}
```

**Typing indicator:**
```json
{ "type": "typing", "chat_id": "uuid" }
```

**Mark read:**
```json
{ "type": "read", "chat_id": "uuid" }
```

**Ping:**
```json
{ "type": "ping" }
```

### Server → Client messages

**New message:**
```json
{
  "type": "message",
  "data": {
    "id": "uuid",
    "chat_id": "uuid",
    "sender": { "id": "uuid", "first_name": "Алексей", "avatar_url": "..." },
    "content": "Привет!",
    "reply_to": null,
    "created_at": "2026-03-15T18:30:00Z",
    "client_id": "temp-uuid"
  }
}
```

**Typing:**
```json
{
  "type": "typing",
  "data": { "chat_id": "uuid", "user_id": "uuid", "first_name": "Алексей" }
}
```

**Read receipt:**
```json
{
  "type": "read",
  "data": { "chat_id": "uuid", "user_id": "uuid", "last_read_at": "..." }
}
```

**Notification (non-chat):**
```json
{
  "type": "notification",
  "data": { "id": "uuid", "type": "result_confirm", "title": "...", "body": "..." }
}
```

**Pong:**
```json
{ "type": "pong" }
```

---

## Summary

| Module | Endpoints | Auth |
|--------|-----------|------|
| Auth | 6 | partial |
| Quiz | 2 | 🔒 |
| Users | 8 | 🔒 |
| Friends | 4 | 🔒 |
| Events | 14 | 🔒 |
| Matches | 5 | 🔒 |
| Communities | 12 | 🔒 |
| Posts | 5 | 🔒 |
| Chat | 7 | 🔒 |
| Notifications | 4 | 🔒 |
| Rating | 4 | 🔒 |
| Courts | 4 | 🔒 |
| Admin | 10 | 🔒 admin+ |
| Superadmin | 7 | 🔒 superadmin |
| WebSocket | 1 | 🔒 |
| **Total** | **~93** | |

---

**Статус:** APPROVED  
**Версия:** 2.0  
**Связанные документы:** TECH-SPEC.md, DATABASE-SCHEMA.sql
