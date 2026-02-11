# Testing Strategy

---

## 1. Overview

| Layer | Type | Tool | Coverage Target |
|-------|------|------|----------------|
| Backend service | Unit tests | Go testing + testify | 80%+ |
| Backend handler | Integration tests | Go httptest | Key flows |
| Backend repository | DB tests | testcontainers-go | Critical queries |
| Mobile | Component tests | Jest + RTL | Shared components |
| Mobile | E2E (post-MVP) | Detox | Critical flows |
| Web | Component tests | Vitest + RTL | Shared components |
| All | Manual testing | Human | Every sprint |

---

## 2. Backend Testing

### Unit Tests (Service Layer) — PRIMARY

Every service method must have a unit test. Use interfaces for dependencies (mock repo, mock SMS, etc).

```go
// internal/service/event_test.go
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type mockEventRepo struct {
    mock.Mock
}

func (m *mockEventRepo) Create(ctx context.Context, input model.CreateEventInput) (*model.Event, error) {
    args := m.Called(ctx, input)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).(*model.Event), args.Error(1)
}

// Table-driven tests
func TestEventService_Create(t *testing.T) {
    tests := []struct {
        name      string
        input     model.CreateEventInput
        mockSetup func(*mockEventRepo)
        wantErr   error
    }{
        {
            name:  "valid event",
            input: model.CreateEventInput{
                Title: "Test Event", EventType: "find_partner",
                LevelMin: 3.0, LevelMax: 4.0, MaxParticipants: 2,
            },
            mockSetup: func(m *mockEventRepo) {
                m.On("Create", mock.Anything, mock.Anything).Return(&model.Event{ID: uuid.New()}, nil)
            },
            wantErr: nil,
        },
        {
            name:  "invalid level range",
            input: model.CreateEventInput{LevelMin: 4.0, LevelMax: 2.0},
            wantErr: ErrInvalidLevelRange,
        },
        {
            name:  "too few participants",
            input: model.CreateEventInput{LevelMin: 3.0, LevelMax: 4.0, MaxParticipants: 1},
            wantErr: ErrTooFewParticipants,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(mockEventRepo)
            if tt.mockSetup != nil { tt.mockSetup(repo) }
            svc := NewEventService(repo)

            _, err := svc.Create(context.Background(), uuid.New(), tt.input)

            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Integration Tests (Handler Layer)

Test full HTTP request → response cycle.

```go
// internal/handler/event_test.go
func TestCreateEvent_Integration(t *testing.T) {
    // Setup
    router := setupTestRouter() // real handler, mock service

    body := `{"title":"Test","event_type":"find_partner","level_min":3.0,"level_max":4.0,"max_participants":2}`
    req := httptest.NewRequest("POST", "/v1/events", strings.NewReader(body))
    req.Header.Set("Authorization", "Bearer "+testToken)
    req.Header.Set("Content-Type", "application/json")
    
    rec := httptest.NewRecorder()
    router.ServeHTTP(rec, req)

    assert.Equal(t, 201, rec.Code)
    
    var resp map[string]any
    json.Unmarshal(rec.Body.Bytes(), &resp)
    assert.NotEmpty(t, resp["id"])
}
```

### ELO Tests (Critical)

```go
// internal/pkg/elo/elo_test.go
func TestCalculate(t *testing.T) {
    tests := []struct {
        name        string
        result      Result
        wantWinDelta float64
        wantLoseDelta float64
    }{
        {
            name: "equal players",
            result: Result{
                WinnerOldRating: 1200, LoserOldRating: 1200,
                WinnerGames: 25, LoserGames: 25,
            },
            wantWinDelta: 16.0, wantLoseDelta: -16.0,
        },
        {
            name: "upset - weak beats strong",
            result: Result{
                WinnerOldRating: 1000, LoserOldRating: 1400,
                WinnerGames: 5, LoserGames: 50,
            },
            wantWinDelta: 36.4, wantLoseDelta: -21.8,
        },
        {
            name: "expected result",
            result: Result{
                WinnerOldRating: 1500, LoserOldRating: 1100,
                WinnerGames: 40, LoserGames: 15,
            },
            wantWinDelta: 2.2, wantLoseDelta: -2.9,
        },
        {
            name: "clamp minimum",
            result: Result{
                WinnerOldRating: 500, LoserOldRating: 110,
                WinnerGames: 50, LoserGames: 50,
            },
            // loser should not go below MinRating (100)
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            change := Calculate(tt.result)
            assert.InDelta(t, tt.wantWinDelta, change.WinnerDelta, 0.5)
            assert.InDelta(t, tt.wantLoseDelta, change.LoserDelta, 0.5)
            assert.GreaterOrEqual(t, change.LoserNewRating, MinRating)
        })
    }
}
```

### What to Test

| Must test (80%+) | Should test | Can skip (MVP) |
|-------------------|------------|----------------|
| Service layer business logic | Handler validation | Repository (sqlc generated) |
| ELO algorithm | Auth middleware | Simple CRUD handlers |
| Auth flow (OTP, JWT) | Rate limiting | Logging |
| Match confirmation logic | WebSocket hub | Config loading |
| Community role checks | File upload validation | |

---

## 3. Frontend Testing (MVP: minimal)

### Shared Components (Jest + React Testing Library)

```typescript
// src/shared/ui/__tests__/Button.test.tsx
import { render, fireEvent } from '@testing-library/react-native';
import { Button } from '../Button';

describe('Button', () => {
  it('renders label', () => {
    const { getByText } = render(<Button label="Click me" onPress={() => {}} />);
    expect(getByText('Click me')).toBeTruthy();
  });

  it('calls onPress', () => {
    const onPress = jest.fn();
    const { getByText } = render(<Button label="Click" onPress={onPress} />);
    fireEvent.press(getByText('Click'));
    expect(onPress).toHaveBeenCalledTimes(1);
  });

  it('disables when loading', () => {
    const onPress = jest.fn();
    const { getByText } = render(<Button label="Click" onPress={onPress} loading />);
    fireEvent.press(getByText('Click'));
    expect(onPress).not.toHaveBeenCalled();
  });
});
```

### API Hooks (if complex logic)

```typescript
// src/features/events/__tests__/useEvents.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useEvents } from '../api';

// Mock API
jest.mock('@/shared/lib/api', () => ({
  api: { get: jest.fn().mockResolvedValue({ data: mockEvents, pagination: mockPagination }) },
}));

it('fetches events', async () => {
  const { result } = renderHook(() => useEvents({}), { wrapper: createWrapper() });
  await waitFor(() => expect(result.current.isSuccess).toBe(true));
  expect(result.current.data).toHaveLength(4);
});
```

---

## 4. Manual Testing Checklist

### Every Sprint: Critical Path

```
□ Auth: Phone → OTP → Profile → Quiz → Home
□ Events: Browse → Detail → Join → Leave
□ Events: Create → Fill → Start → Complete
□ Match: Submit result → Confirm → Rating updated
□ Chat: Open → Send message → Receive message
□ Profile: View → Edit → Save
□ Community: Browse → Join → View leaderboard
□ Notifications: Receive → Tap → Navigate
```

### Pre-Launch: Full Regression

```
□ All auth flows (new user, existing user, re-login, PIN)
□ All CRUD operations (events, communities, posts)
□ All role permissions (member, moderator, admin, owner)
□ Rating calculation (5+ scenarios)
□ Chat (personal, community, event)
□ Push notifications (all types)
□ File upload (avatar, post image)
□ Edge cases (empty states, errors, offline)
□ iOS + Android (both platforms)
□ Localization (ru, kk, en)
□ Deep links
□ App backgrounding and resuming
```

---

## 5. Test Commands

```bash
# Backend
cd apps/backend
go test ./...                          # All tests
go test ./internal/service/... -v      # Service tests verbose
go test ./internal/pkg/elo/... -v      # ELO tests
go test -coverprofile=cover.out ./...  # With coverage
go tool cover -html=cover.out          # View coverage

# Mobile
cd apps/mobile
npm test                               # Jest
npm test -- --coverage                 # With coverage
npm test -- --watch                    # Watch mode

# Web
cd apps/web-admin
npm test                               # Vitest
npm test -- --coverage
```

---

## 6. Test Data / Seed

```go
// scripts/seed/main.go — creates demo data for development and testing

// Users: 20 players with varied ratings (1000-1800)
// Communities: 3 (NTC Astana, Astana Tennis League, Weekend Tennis)
// Events: 10 (mixed types and statuses)
// Matches: 30 (with results and rating history)
// Messages: 50 (in 5 chat rooms)
// Posts: 15 (with likes)
// Badges: assigned to users based on their stats

// Usage:
// make seed
// OR: go run scripts/seed/main.go
```
