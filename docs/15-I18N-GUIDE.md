# i18n Guide — Internationalization
## Локализация теннисной платформы

---

## 1. Языки

| Код | Язык | Статус | Приоритет |
|-----|------|--------|----------|
| `ru` | Русский | Default, полный | Основной |
| `kk` | Казахский | Полный перевод | MVP |
| `en` | English | Полный перевод | MVP |

Язык по умолчанию: `ru`. Определяется по locale устройства, с fallback на `ru`.

---

## 2. Setup

### Mobile (i18next + react-i18next)

```typescript
// src/shared/i18n/index.ts
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import * as Localization from 'expo-localization';

import ru from './locales/ru.json';
import kk from './locales/kk.json';
import en from './locales/en.json';

const deviceLang = Localization.getLocales()[0]?.languageCode || 'ru';

i18n.use(initReactI18next).init({
  resources: { ru: { translation: ru }, kk: { translation: kk }, en: { translation: en } },
  lng: deviceLang === 'kk' ? 'kk' : deviceLang === 'en' ? 'en' : 'ru',
  fallbackLng: 'ru',
  interpolation: { escapeValue: false },
});

export default i18n;
```

### Web Admin (та же библиотека)

```typescript
// Аналогично, но язык из localStorage:
lng: localStorage.getItem('lang') || 'ru',
```

---

## 3. File Structure

```
src/shared/i18n/
├── index.ts              — Setup
├── locales/
│   ├── ru.json           — Русский (основной, пишется первым)
│   ├── kk.json           — Казахский
│   └── en.json           — English
```

---

## 4. Key Naming Convention

### Формат: `{namespace}.{element}`

```
{screen/module}.{component/section}.{element}
```

### Правила

1. Все ключи — lowercase, dot-separated
2. Namespace = экран или модуль (auth, events, profile, chat, common)
3. Без вложенности глубже 3 уровней
4. Для кнопок: `{namespace}.btn_{action}` или `{namespace}.{action}`
5. Для заголовков: `{namespace}.title`
6. Для placeholder: `{namespace}.placeholder_{field}`
7. Для ошибок: `errors.{ERROR_CODE}`

---

## 5. Полная структура ключей

```json
{
  "common": {
    "loading": "Загрузка...",
    "retry": "Повторить",
    "cancel": "Отмена",
    "save": "Сохранить",
    "delete": "Удалить",
    "confirm": "Подтвердить",
    "back": "Назад",
    "done": "Готово",
    "search": "Поиск",
    "search_placeholder": "Поиск...",
    "filter": "Фильтр",
    "sort": "Сортировка",
    "all": "Все",
    "show_more": "Показать ещё",
    "no_results": "Ничего не найдено",
    "pull_to_refresh": "Потяните для обновления",
    "today": "Сегодня",
    "yesterday": "Вчера",
    "just_now": "Только что",
    "minutes_ago": "{{count}} мин. назад",
    "hours_ago": "{{count}} ч. назад",
    "days_ago": "{{count}} дн. назад"
  },

  "tabs": {
    "home": "Главная",
    "players": "Игроки",
    "events": "Ивенты",
    "communities": "Сообщества",
    "profile": "Профиль"
  },

  "auth": {
    "phone_title": "Ваш номер телефона",
    "phone_subtitle": "Мы отправим SMS с кодом подтверждения",
    "phone_placeholder": "+7 (___) ___-__-__",
    "get_code": "Получить код",
    "otp_title": "Введите код",
    "otp_subtitle": "Код отправлен на {{phone}}",
    "resend_code": "Отправить повторно",
    "resend_in": "Повторная отправка через {{seconds}} сек.",
    "wrong_code": "Неверный код",
    "profile_title": "Расскажите о себе",
    "first_name": "Имя",
    "last_name": "Фамилия",
    "gender": "Пол",
    "gender_male": "Мужской",
    "gender_female": "Женский",
    "birth_year": "Год рождения",
    "city": "Город",
    "district": "Район",
    "continue": "Продолжить",
    "quiz_title": "Определим ваш уровень",
    "quiz_skip": "Пропустить",
    "quiz_result_title": "Ваш уровень",
    "quiz_result_subtitle": "NTRP {{level}} — {{name}}",
    "quiz_start": "Начать",
    "pin_set_title": "Установите PIN-код",
    "pin_enter_title": "Введите PIN-код",
    "pin_confirm": "Повторите PIN-код",
    "pin_mismatch": "PIN-коды не совпадают",
    "pin_forgot": "Забыли PIN?",
    "onboarding_1_title": "Находите партнёров",
    "onboarding_1_desc": "Ищите игроков вашего уровня для тренировок и матчей",
    "onboarding_2_title": "Участвуйте в ивентах",
    "onboarding_2_desc": "Турниры, тренировки и парные игры в вашем сообществе",
    "onboarding_3_title": "Следите за рейтингом",
    "onboarding_3_desc": "ELO-рейтинг, статистика матчей и достижения",
    "skip": "Пропустить"
  },

  "home": {
    "title": "Главная",
    "greeting": "Привет, {{name}}!",
    "your_rating": "Ваш рейтинг",
    "upcoming_games": "Ближайшие игры",
    "no_upcoming": "Нет предстоящих игр",
    "quick_find_game": "Найти игру",
    "quick_create_event": "Создать ивент",
    "feed": "Лента",
    "news": "Новости"
  },

  "events": {
    "title": "Ивенты",
    "tab_feed": "Лента",
    "tab_calendar": "Календарь",
    "tab_my": "Мои",
    "create": "Создать ивент",
    "join": "Записаться",
    "leave": "Отписаться",
    "joined": "Вы записаны",
    "spots": "{{current}}/{{max}} мест",
    "spots_left": "Осталось {{count}} мест",
    "full": "Мест нет",
    "status_open": "Открыт",
    "status_filling": "Набор",
    "status_full": "Заполнен",
    "status_in_progress": "Идёт",
    "status_completed": "Завершён",
    "status_cancelled": "Отменён",
    "type_find_partner": "Поиск партнёра",
    "type_organized_game": "Организованная игра",
    "type_tournament": "Турнир",
    "type_training": "Тренировка",
    "composition_singles": "Одиночка",
    "composition_doubles": "Парная",
    "composition_mixed": "Микст",
    "composition_team": "Команды",
    "filter_type": "Тип",
    "filter_level": "Уровень",
    "filter_date": "Дата",
    "filter_district": "Район",
    "detail_participants": "Участники",
    "detail_location": "Место",
    "detail_time": "Время",
    "detail_level": "Уровень {{min}} — {{max}}",
    "detail_sets": "{{count}} сетов",
    "empty": "Нет предстоящих ивентов",
    "empty_my": "Вы ещё не участвовали в ивентах",
    "wizard_step_1": "Тип",
    "wizard_step_2": "Формат",
    "wizard_step_3": "Уровень",
    "wizard_step_4": "Детали",
    "wizard_step_5": "Место",
    "wizard_step_6": "Время",
    "wizard_step_7": "Правила",
    "wizard_step_8": "Проверка"
  },

  "communities": {
    "title": "Сообщества",
    "join": "Вступить",
    "leave": "Выйти",
    "joined": "Вы участник",
    "pending": "Заявка на рассмотрении",
    "members": "{{count}} участников",
    "tab_feed": "Лента",
    "tab_events": "Ивенты",
    "tab_rating": "Рейтинг",
    "tab_members": "Участники",
    "tab_chat": "Чат",
    "type_club": "Клуб",
    "type_league": "Лига",
    "type_organizer": "Организатор",
    "type_group": "Группа",
    "verified": "Верифицировано",
    "create": "Создать сообщество",
    "my": "Мои сообщества",
    "empty": "Нет доступных сообществ",
    "empty_my": "Вы не состоите ни в одном сообществе"
  },

  "players": {
    "title": "Игроки",
    "filter_level": "Уровень",
    "filter_district": "Район",
    "filter_gender": "Пол",
    "rating": "Рейтинг: {{value}}",
    "games": "{{count}} игр",
    "win_rate": "{{value}}% побед",
    "add_friend": "В друзья",
    "remove_friend": "Убрать из друзей",
    "write_message": "Написать",
    "invite_to_game": "Позвать играть",
    "empty": "Нет игроков по вашим фильтрам"
  },

  "profile": {
    "title": "Профиль",
    "edit": "Редактировать",
    "my_stats": "Моя статистика",
    "matches": "Матчи",
    "wins": "Победы",
    "losses": "Поражения",
    "win_rate": "Процент побед",
    "communities": "Сообщества",
    "badges": "Достижения",
    "friends": "Друзья",
    "match_history": "История матчей",
    "rating_chart": "Динамика рейтинга",
    "settings": "Настройки",
    "language": "Язык",
    "notifications_settings": "Уведомления",
    "privacy": "Приватность",
    "about": "О приложении",
    "logout": "Выйти",
    "logout_confirm": "Вы уверены, что хотите выйти?"
  },

  "chat": {
    "title": "Чаты",
    "type_message": "Написать сообщение...",
    "typing": "печатает...",
    "read": "Прочитано",
    "delivered": "Доставлено",
    "empty": "Нет чатов",
    "mute": "Отключить уведомления",
    "unmute": "Включить уведомления"
  },

  "notifications": {
    "title": "Уведомления",
    "mark_all_read": "Прочитать все",
    "empty": "Нет уведомлений",
    "match_result_pending": "Подтвердите результат матча",
    "match_result_confirmed": "Результат подтверждён",
    "rating_changed": "Рейтинг обновлён: {{rating}} ({{delta}})",
    "event_reminder": "{{event}} начнётся через 1 час",
    "event_joined": "{{user}} записался на {{event}}",
    "new_message": "{{sender}}: {{preview}}",
    "badge_earned": "Новое достижение: {{badge}}"
  },

  "match": {
    "submit_result": "Внести результат",
    "confirm_result": "Подтвердить результат",
    "dispute_result": "Оспорить",
    "score": "Счёт",
    "set": "Сет {{number}}",
    "tiebreak": "Тай-брейк",
    "winner": "Победитель",
    "result_submitted": "Результат отправлен, ожидаем подтверждения",
    "result_confirmed": "Результат подтверждён!"
  },

  "errors": {
    "something_went_wrong": "Что-то пошло не так. Попробуйте позже",
    "no_internet": "Нет подключения к интернету",
    "timeout": "Сервер не отвечает",
    "UNAUTHORIZED": "Требуется авторизация",
    "TOKEN_EXPIRED": "Сессия истекла",
    "OTP_SESSION_EXPIRED": "Код истёк, запросите новый",
    "OTP_INVALID_CODE": "Неверный код",
    "OTP_MAX_ATTEMPTS": "Слишком много попыток",
    "EVENT_FULL": "Все места заняты",
    "EVENT_CLOSED": "Запись закрыта",
    "EVENT_WRONG_LEVEL": "Ваш уровень не подходит",
    "ALREADY_MEMBER": "Вы уже участник",
    "RATE_LIMITED": "Подождите немного",
    "SMS_RATE_LIMITED": "Подождите перед повторной отправкой"
  }
}
```

---

## 6. Workflow для добавления нового текста

1. Добавь ключ в `ru.json` (основной файл)
2. Добавь перевод в `kk.json` и `en.json`
3. Используй в компоненте: `const { t } = useTranslation(); t('events.join')`
4. НИКОГДА не пиши текст напрямую: `<Text>Записаться</Text>` ❌
5. ВСЕГДА через t(): `<Text>{t('events.join')}</Text>` ✅

---

## 7. Pluralization

```json
{
  "events": {
    "spots_left_one": "Осталось {{count}} место",
    "spots_left_few": "Осталось {{count}} места",
    "spots_left_many": "Осталось {{count}} мест"
  }
}
```

```typescript
t('events.spots_left', { count: 3 }) // → "Осталось 3 места"
```

---

## 8. Date/Time Formatting

```typescript
// Используй Intl.DateTimeFormat, не hardcode форматы
const formatDate = (date: string, locale: string) =>
  new Intl.DateTimeFormat(locale, { day: 'numeric', month: 'short' }).format(new Date(date));

// ru: "15 фев."  kk: "15 ақп."  en: "Feb 15"
```

---

## 9. Казахский язык — особенности

- Писать на кириллице (не латинице, переход ещё не завершён)
- Pluralization: правила отличаются от русского
- Некоторые термины не переводятся: NTRP, ELO, PIN, OTP
- Спортивные термины: использовать устоявшиеся казахские варианты где есть
