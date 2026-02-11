# Analytics Events Specification
## Трекинг событий для Firebase Analytics

---

## 1. Setup

**Инструмент:** Firebase Analytics (free, unlimited events)
**Mobile:** `expo-firebase-analytics` / `@react-native-firebase/analytics`
**Web:** Firebase Analytics web SDK

---

## 2. Event Naming

Convention: `snake_case`, verb_noun format, max 40 chars.

---

## 3. Events Catalog

### Auth

| Event | Когда | Parameters |
|-------|-------|-----------|
| `app_open` | Каждый запуск | `is_authenticated`, `app_version` |
| `onboarding_start` | Показан первый onboarding экран | — |
| `onboarding_complete` | Нажал "Начать" на последнем | — |
| `onboarding_skip` | Нажал "Пропустить" | `screen_index` |
| `auth_phone_submit` | Отправил номер | — |
| `auth_otp_submit` | Ввёл OTP код | `is_valid` |
| `auth_otp_resend` | Нажал "Отправить повторно" | — |
| `auth_profile_complete` | Заполнил профиль | `gender`, `district` |
| `auth_quiz_start` | Начал квиз | — |
| `auth_quiz_complete` | Завершил квиз | `ntrp_level` |
| `auth_quiz_skip` | Пропустил квиз | — |
| `auth_pin_set` | Установил PIN | — |
| `auth_login` | Успешный вход (OTP или PIN) | `method`: otp/pin |
| `auth_logout` | Выход из аккаунта | — |

### Navigation

| Event | Когда | Parameters |
|-------|-------|-----------|
| `screen_view` | Открытие любого экрана | `screen_name`, `screen_class` |
| `tab_switch` | Переключение таба | `tab_name` |

### Events (Ивенты)

| Event | Когда | Parameters |
|-------|-------|-----------|
| `event_view` | Открыл детали ивента | `event_id`, `event_type` |
| `event_create_start` | Начал wizard создания | — |
| `event_create_complete` | Создал ивент | `event_type`, `composition_type`, `max_participants` |
| `event_create_abandon` | Закрыл wizard без создания | `step` (на каком шаге) |
| `event_join` | Записался на ивент | `event_id`, `event_type` |
| `event_leave` | Отписался от ивента | `event_id` |
| `event_filter` | Применил фильтр | `filter_type`, `filter_value` |
| `event_search` | Использовал поиск | `query_length` |

### Matches

| Event | Когда | Parameters |
|-------|-------|-----------|
| `match_result_submit` | Внёс результат | `match_id`, `sets_count` |
| `match_result_confirm` | Подтвердил результат | `match_id` |
| `match_result_dispute` | Оспорил результат | `match_id` |

### Communities

| Event | Когда | Parameters |
|-------|-------|-----------|
| `community_view` | Открыл сообщество | `community_id`, `community_type` |
| `community_join` | Вступил | `community_id`, `community_type` |
| `community_leave` | Вышел | `community_id` |
| `community_create` | Создал сообщество | `community_type` |

### Social

| Event | Когда | Parameters |
|-------|-------|-----------|
| `chat_open` | Открыл чат | `chat_type`: personal/community/event |
| `chat_message_send` | Отправил сообщение | `chat_type` |
| `post_create` | Создал пост | `has_image` |
| `post_like` | Лайкнул пост | — |
| `friend_add` | Добавил в друзья | — |
| `player_profile_view` | Открыл чужой профиль | — |

### Settings

| Event | Когда | Parameters |
|-------|-------|-----------|
| `language_change` | Сменил язык | `language`: ru/kk/en |
| `notification_toggle` | Вкл/выкл уведомления | `type`, `enabled` |

---

## 4. User Properties

Устанавливаются один раз, обновляются при изменении:

| Property | Тип | Пример |
|----------|-----|--------|
| `ntrp_level` | string | "3.0" |
| `rating` | number | 1200 |
| `gender` | string | "male" |
| `district` | string | "Есильский" |
| `communities_count` | number | 3 |
| `games_played` | number | 25 |
| `app_language` | string | "ru" |

---

## 5. Ключевые метрики для отслеживания

| Метрика | Формула | Цель (6 мес.) |
|---------|---------|---------------|
| DAU | Уникальные пользователи/день | 50+ |
| WAU | Уникальные пользователи/неделя | 150+ |
| MAU | Уникальные пользователи/месяц | 300+ |
| Retention D1 | % вернувшихся на следующий день | >40% |
| Retention D7 | % вернувшихся через неделю | >25% |
| Retention D30 | % вернувшихся через месяц | >15% |
| Auth completion | onboarding → auth → profile → quiz → home | >60% |
| Event join rate | event_view → event_join | >20% |
| Match confirm rate | result_submit → result_confirm | >80% |

---

## 6. Implementation

```typescript
// src/shared/lib/analytics.ts
import * as Analytics from 'expo-firebase-analytics';

export const analytics = {
  track: (event: string, params?: Record<string, any>) => {
    if (__DEV__) {
      console.log('[Analytics]', event, params);
      return;
    }
    Analytics.logEvent(event, params);
  },

  setUser: (userId: string) => Analytics.setUserId(userId),
  setProperty: (name: string, value: string) => Analytics.setUserProperty(name, value),
  screenView: (name: string) => Analytics.logEvent('screen_view', { screen_name: name }),
};

// Usage:
analytics.track('event_join', { event_id: '...', event_type: 'tournament' });
analytics.setProperty('ntrp_level', '3.0');
```
