# Error Codes Catalog
## Единый каталог кодов ошибок

---

## Response Format

Все ошибки возвращаются в формате:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Описание на языке пользователя",
    "details": [
      { "field": "phone", "message": "Неверный формат номера" }
    ]
  }
}
```

`details` — опционально, только для VALIDATION_ERROR.

---

## HTTP Status → Error Code Mapping

| HTTP | Категория | Когда |
|------|----------|-------|
| 400 | Client error | Невалидный input, бизнес-правило нарушено |
| 401 | Auth error | Нет токена, токен невалидный, истёк |
| 403 | Permission error | Нет прав на действие |
| 404 | Not found | Ресурс не найден |
| 409 | Conflict | Дублирование, уже существует |
| 429 | Rate limit | Превышен лимит запросов |
| 500 | Server error | Внутренняя ошибка (никогда не раскрывать детали клиенту) |

---

## Полный каталог

### Auth Errors (401)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| UNAUTHORIZED | Требуется авторизация | Authorization required | Нет Authorization header |
| TOKEN_EXPIRED | Токен истёк | Token expired | JWT access token expired |
| TOKEN_INVALID | Невалидный токен | Invalid token | JWT подпись не проходит, malformed |
| TOKEN_REVOKED | Токен отозван | Token revoked | Refresh token уже использован (possible compromise) |

### OTP Errors (400)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| OTP_SESSION_EXPIRED | Сессия истекла, запросите код повторно | Session expired, request a new code | Redis session TTL expired или не найдена |
| OTP_INVALID_CODE | Неверный код | Invalid code | Код не совпадает |
| OTP_MAX_ATTEMPTS | Превышено количество попыток | Too many attempts | 5 неверных попыток за сессию |

### PIN Errors (400)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| PIN_INVALID | Неверный PIN-код | Invalid PIN | PIN не совпадает с hash |
| PIN_LOCKED | PIN заблокирован, войдите через SMS | PIN locked, use SMS login | 3 неверных попытки |
| PIN_ALREADY_SET | PIN уже установлен | PIN already set | Повторная установка |

### Validation Errors (400)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| VALIDATION_ERROR | Ошибка валидации | Validation error | Любая ошибка валидации полей (details содержит конкретные поля) |
| INVALID_PHONE_FORMAT | Неверный формат номера | Invalid phone format | Телефон не проходит regex `^\+7[0-9]{10}$` |
| INVALID_JSON | Невалидный JSON | Invalid JSON | Body не парсится как JSON |

### Permission Errors (403)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| FORBIDDEN | Недостаточно прав | Access denied | Общая ошибка прав доступа |
| NOT_COMMUNITY_MEMBER | Вы не участник сообщества | Not a community member | Действие требует членства |
| INSUFFICIENT_ROLE | Недостаточный уровень роли | Insufficient role | Действие требует admin/moderator/owner |

### Not Found (404)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| NOT_FOUND | Не найдено | Not found | Общий 404 |
| USER_NOT_FOUND | Пользователь не найден | User not found | GET /users/:id — нет такого |
| EVENT_NOT_FOUND | Ивент не найден | Event not found | GET /events/:id — нет или deleted |
| COMMUNITY_NOT_FOUND | Сообщество не найдено | Community not found | |
| MATCH_NOT_FOUND | Матч не найден | Match not found | |
| CHAT_NOT_FOUND | Чат не найден | Chat not found | |

### Conflict (409)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| ALREADY_EXISTS | Уже существует | Already exists | Общий конфликт |
| ALREADY_MEMBER | Вы уже участник | Already a member | POST /communities/:id/join — уже в сообществе |
| ALREADY_JOINED_EVENT | Вы уже записаны | Already joined | POST /events/:id/join — уже записан |
| ALREADY_FRIENDS | Уже в друзьях | Already friends | POST /friends — уже друзья |
| PROFILE_ALREADY_SET | Профиль уже заполнен | Profile already set | POST /auth/profile/setup — повторно |
| RESULT_ALREADY_SUBMITTED | Результат уже внесён | Result already submitted | POST /matches/:id/result — повторно |

### Business Logic Errors (400)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| EVENT_FULL | Все места заняты | Event is full | POST /events/:id/join — max_participants reached |
| EVENT_CLOSED | Ивент закрыт для записи | Event registration closed | Статус не open/filling |
| EVENT_WRONG_LEVEL | Ваш уровень не подходит | Level requirement not met | Уровень игрока вне [level_min, level_max] |
| EVENT_INVALID_STATUS_TRANSITION | Недопустимая смена статуса | Invalid status transition | open→completed (без in_progress) |
| MATCH_NOT_YOUR_TURN | Ожидание подтверждения соперника | Waiting for opponent confirmation | Пытается confirm свой же результат |
| MATCH_ALREADY_CONFIRMED | Матч уже подтверждён | Match already confirmed | Повторное подтверждение |
| MATCH_DISPUTED | Результат оспорен | Result disputed | Второй игрок не согласен |
| COMMUNITY_REQUIRES_VERIFICATION | Требуется верификация | Verification required | Для клубов/лиг нужна верификация суперадмином |
| SELF_ACTION_NOT_ALLOWED | Действие с самим собой невозможно | Cannot perform action on yourself | Добавить себя в друзья, играть против себя |
| CANNOT_LEAVE_AS_OWNER | Передайте владение перед выходом | Transfer ownership before leaving | Owner пытается выйти из сообщества |

### Rate Limit (429)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| RATE_LIMITED | Слишком много запросов | Too many requests | Общий rate limit |
| SMS_RATE_LIMITED | Подождите перед повторной отправкой | Wait before resending | SMS 3/час или 10/день |

### Server Errors (500)

| Code | Message (ru) | Message (en) | Когда |
|------|-------------|--------------|-------|
| INTERNAL_ERROR | Внутренняя ошибка сервера | Internal server error | Любая необработанная ошибка. НИКОГДА не показывать stack trace клиенту |

---

## Frontend Handling

### Mobile (React Native)
```typescript
// src/shared/lib/errors.ts

type ErrorCode = 
  | 'UNAUTHORIZED' | 'TOKEN_EXPIRED' | 'TOKEN_INVALID' | 'TOKEN_REVOKED'
  | 'OTP_SESSION_EXPIRED' | 'OTP_INVALID_CODE' | 'OTP_MAX_ATTEMPTS'
  | 'VALIDATION_ERROR' | 'INVALID_PHONE_FORMAT'
  | 'FORBIDDEN' | 'NOT_COMMUNITY_MEMBER' | 'INSUFFICIENT_ROLE'
  | 'NOT_FOUND' | 'USER_NOT_FOUND' | 'EVENT_NOT_FOUND'
  | 'ALREADY_MEMBER' | 'ALREADY_JOINED_EVENT' | 'EVENT_FULL'
  | 'EVENT_CLOSED' | 'EVENT_WRONG_LEVEL' | 'RATE_LIMITED'
  | 'INTERNAL_ERROR';

function handleApiError(error: ApiError) {
  switch (error.code) {
    // Auth errors → redirect to login
    case 'UNAUTHORIZED':
    case 'TOKEN_EXPIRED':
    case 'TOKEN_REVOKED':
      authStore.logout();
      break;

    // Validation → show field errors
    case 'VALIDATION_ERROR':
      return error.details; // for form.setErrors()

    // Business logic → show toast
    case 'EVENT_FULL':
    case 'EVENT_CLOSED':
    case 'ALREADY_MEMBER':
      showToast(error.message);
      break;

    // Rate limit → show toast with retry
    case 'RATE_LIMITED':
    case 'SMS_RATE_LIMITED':
      showToast(error.message);
      break;

    // Server error → generic message
    case 'INTERNAL_ERROR':
    default:
      showToast(t('errors.something_went_wrong'));
  }
}
```

### i18n Error Keys

```json
{
  "errors": {
    "something_went_wrong": "Что-то пошло не так. Попробуйте позже",
    "no_internet": "Нет подключения к интернету",
    "timeout": "Сервер не отвечает",
    "UNAUTHORIZED": "Требуется авторизация",
    "TOKEN_EXPIRED": "Сессия истекла, войдите заново",
    "OTP_SESSION_EXPIRED": "Код истёк, запросите новый",
    "OTP_INVALID_CODE": "Неверный код",
    "OTP_MAX_ATTEMPTS": "Слишком много попыток",
    "EVENT_FULL": "Все места заняты",
    "EVENT_CLOSED": "Запись закрыта",
    "EVENT_WRONG_LEVEL": "Ваш уровень не подходит",
    "ALREADY_MEMBER": "Вы уже участник",
    "RATE_LIMITED": "Подождите немного"
  }
}
```

---

## Backend Implementation

```go
// internal/service/errors.go

var (
    // Auth
    ErrUnauthorized    = &AppError{Code: "UNAUTHORIZED", Status: 401}
    ErrTokenExpired    = &AppError{Code: "TOKEN_EXPIRED", Status: 401}
    ErrTokenInvalid    = &AppError{Code: "TOKEN_INVALID", Status: 401}
    ErrTokenRevoked    = &AppError{Code: "TOKEN_REVOKED", Status: 401}

    // OTP
    ErrOTPSessionExpired = &AppError{Code: "OTP_SESSION_EXPIRED", Status: 400}
    ErrOTPInvalidCode    = &AppError{Code: "OTP_INVALID_CODE", Status: 400}
    ErrOTPMaxAttempts    = &AppError{Code: "OTP_MAX_ATTEMPTS", Status: 400}

    // Permission
    ErrForbidden         = &AppError{Code: "FORBIDDEN", Status: 403}
    ErrNotCommunityMember = &AppError{Code: "NOT_COMMUNITY_MEMBER", Status: 403}
    ErrInsufficientRole  = &AppError{Code: "INSUFFICIENT_ROLE", Status: 403}

    // Not Found
    ErrNotFound          = &AppError{Code: "NOT_FOUND", Status: 404}
    ErrUserNotFound      = &AppError{Code: "USER_NOT_FOUND", Status: 404}
    ErrEventNotFound     = &AppError{Code: "EVENT_NOT_FOUND", Status: 404}

    // Conflict
    ErrAlreadyMember      = &AppError{Code: "ALREADY_MEMBER", Status: 409}
    ErrAlreadyJoinedEvent = &AppError{Code: "ALREADY_JOINED_EVENT", Status: 409}
    ErrEventFull          = &AppError{Code: "EVENT_FULL", Status: 400}
    ErrEventClosed        = &AppError{Code: "EVENT_CLOSED", Status: 400}

    // Rate Limit
    ErrRateLimited       = &AppError{Code: "RATE_LIMITED", Status: 429}
    ErrSMSRateLimited    = &AppError{Code: "SMS_RATE_LIMITED", Status: 429}
)

type AppError struct {
    Code    string `json:"code"`
    Status  int    `json:"-"`
    Message string `json:"message,omitempty"`
}

func (e *AppError) Error() string { return e.Code }

func (e *AppError) WithMessage(msg string) *AppError {
    return &AppError{Code: e.Code, Status: e.Status, Message: msg}
}
```
