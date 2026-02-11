# Integrations Specification
## Внешние сервисы и SDK

---

## 1. SMS Provider — smsc.kz

### Purpose
OTP codes for phone authentication.

### Config
```env
SMS_PROVIDER=smsc          # smsc | mobizon | mock
SMS_API_URL=https://smsc.kz/sys/send.php
SMS_API_LOGIN=your_login
SMS_API_PASSWORD=your_password
SMS_SENDER=TennisApp
```

### API Call
```go
// internal/pkg/sms/smsc.go

func (s *SMSCProvider) Send(phone, message string) error {
    params := url.Values{
        "login":  {s.config.Login},
        "psw":    {s.config.Password},
        "phones": {phone},           // format: 77071234567
        "mes":    {message},
        "sender": {s.config.Sender},
        "charset": {"utf-8"},
        "fmt":    {"3"},              // JSON response
    }
    resp, err := http.Get(s.config.URL + "?" + params.Encode())
    // Handle response...
}
```

### OTP Message Format
```
Ваш код: 1234

Tennis Astana
```

### Mock Provider (development)
```go
// internal/pkg/sms/mock.go
// В dev режиме: код всегда "1234", SMS не отправляется, логируется в slog
func (m *MockProvider) Send(phone, message string) error {
    slog.Info("mock SMS", "phone", phone, "message", message)
    return nil
}
```

### Rate Limits
- 3 SMS / hour per phone number
- 10 SMS / day per phone number
- OTP valid for 5 minutes
- Max 5 verification attempts per session

---

## 2. Firebase Cloud Messaging (FCM)

### Purpose
Push notifications for mobile app.

### Config
```env
FIREBASE_CREDENTIALS=/path/to/service-account.json
# OR
FIREBASE_PROJECT_ID=tennis-astana
FIREBASE_PRIVATE_KEY=...
FIREBASE_CLIENT_EMAIL=...
```

### Setup
```go
// internal/pkg/firebase/firebase.go
import (
    firebase "firebase.google.com/go/v4"
    "firebase.google.com/go/v4/messaging"
)

func NewClient(credentialsPath string) (*messaging.Client, error) {
    opt := option.WithCredentialsFile(credentialsPath)
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil { return nil, err }
    return app.Messaging(context.Background())
}
```

### Send Notification
```go
func (f *FCMService) SendToUser(ctx context.Context, userID uuid.UUID, title, body string, data map[string]string) error {
    // 1. Get user's FCM tokens from DB
    tokens, err := f.repo.GetUserFCMTokens(ctx, userID)
    if err != nil { return err }

    // 2. Send to each token
    for _, token := range tokens {
        msg := &messaging.Message{
            Token: token,
            Notification: &messaging.Notification{
                Title: title,
                Body:  body,
            },
            Data: data,
            Android: &messaging.AndroidConfig{
                Priority: "high",
            },
            APNS: &messaging.APNSConfig{
                Payload: &messaging.APNSPayload{
                    Aps: &messaging.Aps{
                        Sound: "default",
                        Badge: &unreadCount,
                    },
                },
            },
        }
        _, err := f.client.Send(ctx, msg)
        if messaging.IsUnregistered(err) {
            // Token expired, remove from DB
            f.repo.DeleteFCMToken(ctx, token)
        }
    }
    return nil
}
```

### Notification Types
```go
var NotificationTemplates = map[string]struct{ Title, Body string }{
    "match_result_pending":  {"Подтвердите результат", "{{opponent}} внёс результат: {{score}}"},
    "match_result_confirmed": {"Результат подтверждён", "{{score}} — рейтинг обновлён"},
    "rating_changed":         {"Рейтинг обновлён", "Ваш рейтинг: {{rating}} ({{delta}})"},
    "event_reminder":         {"Напоминание", "{{event}} начнётся через 1 час"},
    "event_joined":           {"Новый участник", "{{user}} записался на {{event}}"},
    "new_message":             {"Новое сообщение", "{{sender}}: {{preview}}"},
    "community_verified":      {"Сообщество подтверждено", "{{community}} прошло верификацию"},
    "badge_earned":            {"Новое достижение!", "Вы получили бейдж «{{badge}}»"},
}
```

### Mobile Registration
```typescript
// apps/mobile/src/shared/lib/notifications.ts
import * as Notifications from 'expo-notifications';
import { api } from './api';

export async function registerForPush() {
  const { status } = await Notifications.requestPermissionsAsync();
  if (status !== 'granted') return;
  
  const token = (await Notifications.getExpoPushTokenAsync()).data;
  // OR for FCM directly:
  // const token = (await Notifications.getDevicePushTokenAsync()).data;
  
  await api.post('/v1/users/me/fcm-token', { token, platform: Platform.OS });
}
```

---

## 3. Object Storage — Cloudflare R2 / MinIO

### Purpose
User avatars, community logos, post images.

### Config
```env
S3_ENDPOINT=http://localhost:9000         # MinIO local
S3_ENDPOINT=https://xxx.r2.cloudflarestorage.com  # R2 production
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=tennisapp
S3_REGION=auto
S3_PUBLIC_URL=https://cdn.tennisapp.kz    # CDN URL for public access
```

### Implementation
```go
// internal/pkg/storage/s3.go
import (
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
    client    *s3.Client
    bucket    string
    publicURL string
}

func (s *Storage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
    _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      &s.bucket,
        Key:         &key,
        Body:        reader,
        ContentType: &contentType,
    })
    if err != nil { return "", err }
    return fmt.Sprintf("%s/%s", s.publicURL, key), nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
    _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: &s.bucket,
        Key:    &key,
    })
    return err
}
```

### File Organization
```
tennisapp/
├── avatars/{user_id}/avatar.jpg
├── communities/{community_id}/logo.jpg
├── posts/{post_id}/{filename}.jpg
└── courts/{court_id}/photo.jpg
```

### Upload Flow
```
1. Mobile: pick image → compress → POST /v1/upload (multipart)
2. Backend: validate (type, size) → generate key → upload to S3 → return URL
3. Backend: save URL to relevant table (users.avatar_url, etc.)
```

### Limits
- Max file size: 5 MB
- Allowed types: image/jpeg, image/png, image/webp
- Max dimensions: 2048x2048 (resize server-side in Phase 2)

---

## 4. Google Maps Platform

### Purpose
Courts map on mobile app.

### Config
```env
GOOGLE_MAPS_API_KEY=AIza...  # Restricted to Maps SDK for iOS/Android
```

### Mobile Usage
```typescript
// react-native-maps (Expo compatible)
import MapView, { Marker } from 'react-native-maps';

<MapView
  region={{
    latitude: 51.1694,  // Astana center
    longitude: 71.4491,
    latitudeDelta: 0.1,
    longitudeDelta: 0.1,
  }}
>
  {courts.map(court => (
    <Marker
      key={court.id}
      coordinate={{ latitude: court.latitude, longitude: court.longitude }}
      title={court.name}
      description={court.address}
    />
  ))}
</MapView>
```

### Courts Data
Courts are stored in PostgreSQL `courts` table with lat/lng. No geocoding API needed — admin enters coordinates manually via superadmin panel.

---

## 5. Sentry — Error Tracking

### Config
```env
SENTRY_DSN=https://xxx@sentry.io/yyy
SENTRY_ENVIRONMENT=production
```

### Backend
```go
import "github.com/getsentry/sentry-go"

func init() {
    sentry.Init(sentry.ClientOptions{
        Dsn:         os.Getenv("SENTRY_DSN"),
        Environment: os.Getenv("SENTRY_ENVIRONMENT"),
    })
}
```

### Mobile
```typescript
import * as Sentry from '@sentry/react-native';

Sentry.init({
  dsn: process.env.EXPO_PUBLIC_SENTRY_DSN,
  environment: __DEV__ ? 'development' : 'production',
});
```

---

## 6. Redis — Cache & Pub/Sub

### Purpose
JWT refresh token storage, rate limiting, WebSocket pub/sub, cache.

### Config
```env
REDIS_URL=redis://localhost:6379           # Local
REDIS_URL=rediss://xxx.upstash.io:6379     # Upstash (production)
```

### Key Patterns
```
refresh:{token_hash}          → user_id (TTL: 30 days)
otp:{session_id}              → {phone, code, attempts} (TTL: 5 min)
rate:sms:{phone}:hour         → count (TTL: 1 hour)
rate:sms:{phone}:day          → count (TTL: 24 hours)
rate:api:{user_id}            → count (TTL: 1 min)
ws:channel:{chat_id}          → pub/sub channel for WebSocket
cache:user:{user_id}          → user profile JSON (TTL: 5 min)
cache:leaderboard:{community} → leaderboard JSON (TTL: 1 min)
online:{user_id}              → "1" (TTL: 5 min, refreshed on activity)
```

---

## 7. Integration Priority

| Sprint | Integration | Effort |
|--------|-------------|--------|
| 1 | Redis (cache, rate limit) | 2h |
| 1 | SMS mock provider | 1h |
| 2 | SMS real provider (smsc.kz) | 2h |
| 5 | Firebase FCM | 4h |
| 5 | WebSocket + Redis pub/sub | 8h |
| 7 | MinIO / S3 (file upload) | 4h |
| 7 | Google Maps (courts) | 2h |
| 11 | Sentry (error tracking) | 1h |
