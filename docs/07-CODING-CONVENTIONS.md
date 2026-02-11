# Coding Conventions
## Стандарты кода теннисной платформы

---

## 1. Go Backend

### Структура файла
```go
package service

import (
    "context"        // stdlib first
    "fmt"

    "github.com/google/uuid"  // third-party second
    
    "tennisapp/internal/model" // project imports last
    "tennisapp/internal/repository"
)
```

### Naming

| Элемент | Стиль | Пример |
|---------|-------|--------|
| Package | lowercase single word | `handler`, `service`, `repository` |
| Struct | PascalCase | `EventService`, `CreateEventRequest` |
| Interface | PascalCase, -er suffix | `EventCreator`, `UserRepository` |
| Method (exported) | PascalCase | `func (s *Service) Create()` |
| Method (private) | camelCase | `func (s *Service) validateLevel()` |
| Variable | camelCase | `eventCount`, `userID` |
| Constant | PascalCase | `MaxParticipants`, `DefaultPageSize` |
| Env variable | SCREAMING_SNAKE | `DATABASE_URL`, `JWT_SECRET` |
| DB column | snake_case | `created_at`, `user_id` |
| JSON field | snake_case | `"first_name"`, `"win_rate"` |
| URL path | kebab-case | `/match-results`, `/otp/send` |
| Error var | Err prefix | `ErrNotFound`, `ErrForbidden` |

### Error Handling

```go
// 1. Define typed errors in service
var (
    ErrNotFound      = errors.New("not found")
    ErrForbidden     = errors.New("forbidden")
    ErrAlreadyExists = errors.New("already exists")
    ErrInvalidInput  = errors.New("invalid input")
)

// 2. Wrap with context
return nil, fmt.Errorf("create event: %w", err)

// 3. NEVER ignore errors
result, err := s.repo.Get(ctx, id)
if err != nil {
    return nil, fmt.Errorf("get event %s: %w", id, err)
}

// 4. Handler maps to HTTP
switch {
case errors.Is(err, service.ErrNotFound):    // → 404
case errors.Is(err, service.ErrForbidden):   // → 403
case errors.Is(err, service.ErrInvalidInput): // → 400
default:                                      // → 500 + log
}
```

### Context Propagation
```go
// ALWAYS pass context through layers
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    result, err := h.service.Get(ctx, id)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*model.Event, error) {
    return s.repo.GetByID(ctx, id)
}
```

### Validation
```go
// Use struct tags for request validation
type CreateEventRequest struct {
    Title           string  `json:"title" validate:"required,min=3,max=100"`
    EventType       string  `json:"event_type" validate:"required,oneof=find_partner organized_game tournament training"`
    CompositionType string  `json:"composition_type" validate:"required,oneof=singles doubles mixed team"`
    LevelMin        float64 `json:"level_min" validate:"required,gte=1.0,lte=7.0"`
    LevelMax        float64 `json:"level_max" validate:"required,gte=1.0,lte=7.0,gtefield=LevelMin"`
    MaxParticipants int     `json:"max_participants" validate:"required,gte=2,lte=64"`
}
```

### SQL (sqlc)
```sql
-- File: repository/queries/event.sql
-- ALWAYS name queries with pattern: Verb + Entity
-- ALWAYS specify return type: :one, :many, :exec, :execrows

-- name: GetEventByID :one
SELECT id, community_id, creator_id, title, event_type, status,
       level_min, level_max, max_participants, start_time, created_at
FROM events
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEventsByCommunity :many
SELECT id, title, event_type, status, start_time, max_participants,
       (SELECT COUNT(*) FROM event_participants WHERE event_id = events.id) as participant_count
FROM events
WHERE community_id = $1
  AND deleted_at IS NULL
  AND ($2::text = '' OR status = $2)
ORDER BY start_time DESC
LIMIT $3 OFFSET $4;

-- NEVER: SELECT * (always list columns)
-- ALWAYS: WHERE deleted_at IS NULL
-- ALWAYS: LIMIT + OFFSET for lists
-- ALWAYS: ORDER BY for deterministic results
```

---

## 2. TypeScript (Mobile + Web)

### Naming

| Элемент | Стиль | Пример |
|---------|-------|--------|
| Component file | PascalCase.tsx | `EventCard.tsx` |
| Hook file | camelCase.ts | `useEvents.ts` |
| Utility file | camelCase.ts | `formatDate.ts` |
| Type/Interface | PascalCase | `Event`, `CreateEventInput` |
| Function | camelCase | `formatDate`, `calculateRating` |
| Variable | camelCase | `eventCount`, `isLoading` |
| Constant | SCREAMING_SNAKE | `API_BASE_URL`, `MAX_RETRIES` |
| Component | PascalCase | `EventCard`, `ProfileScreen` |
| Hook | camelCase, use- prefix | `useEvents`, `useAuth` |
| Event handler | handle- prefix | `handlePress`, `handleSubmit` |
| Boolean | is/has/can prefix | `isLoading`, `hasError`, `canEdit` |

### Types

```typescript
// ALWAYS define types, NEVER use `any`
interface Event {
  id: string;
  title: string;
  event_type: EventType;
  status: EventStatus;
  level_min: number;
  level_max: number;
  max_participants: number;
  participant_count: number;
  start_time: string; // ISO 8601
  created_at: string;
}

type EventType = 'find_partner' | 'organized_game' | 'tournament' | 'training';
type EventStatus = 'draft' | 'open' | 'filling' | 'full' | 'in_progress' | 'completed' | 'cancelled';

// API response types
interface ApiResponse<T> {
  data: T;
  pagination?: Pagination;
}

interface ApiError {
  error: {
    code: string;
    message: string;
    details?: { field: string; message: string }[];
  };
}
```

### Component Rules

```tsx
// 1. Props interface ALWAYS defined
interface EventCardProps {
  event: Event;
  onPress: () => void;
  variant?: 'default' | 'compact';
}

// 2. Named exports (not default) for components
export function EventCard({ event, onPress, variant = 'default' }: EventCardProps) {
  // 3. Hooks at top
  const { t } = useTranslation();
  const [isExpanded, setIsExpanded] = useState(false);

  // 4. Early returns for edge cases
  if (!event) return null;

  // 5. Derived values
  const spotsLeft = event.max_participants - event.participant_count;
  const isFull = spotsLeft <= 0;

  return (/* JSX */);
}

// 6. Styles at bottom (mobile)
const styles = StyleSheet.create({
  container: { padding: 16, borderRadius: 16 },
});
```

### File Size Limits
- Component: max 200 lines → split into subcomponents
- Hook: max 100 lines → extract helpers
- Utility: max 50 lines per function

---

## 3. Universal Rules

### Git Commits
```
<type>(<scope>): <description>

type:  feat | fix | refactor | test | docs | chore | style
scope: auth | events | chat | users | communities | rating | admin | mobile | web | ci
```

### Comments
```go
// Comment full sentences. Explain WHY, not WHAT.

// calculateELO applies the modified ELO algorithm with variable K-factor
// based on the player's total number of games. New players (< 10 games)
// have a higher K-factor for faster initial calibration.
func calculateELO(winner, loser float64, games int) (float64, float64) {
```

```typescript
// Comment complex business logic only.
// Self-documenting code > comments for obvious things.

// WRONG:
// Set loading to true
setIsLoading(true);

// RIGHT:
// Debounce search to avoid excessive API calls during fast typing
const debouncedSearch = useDebounce(searchQuery, 300);
```

### No Magic Numbers
```go
// WRONG
if attempts > 5 { ... }

// RIGHT
const MaxOTPAttempts = 5
if attempts > MaxOTPAttempts { ... }
```

```typescript
// WRONG
{ marginTop: 16, borderRadius: 12 }

// RIGHT
import { spacing, radius } from '@/shared/theme';
{ marginTop: spacing.md, borderRadius: radius.md }
```
