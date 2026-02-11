# Security Specification

---

## 1. Authentication Flow

```
Phone Input → OTP Send → OTP Verify → [Profile Setup] → [Quiz] → JWT Issued
                                                                      │
                                                           ┌──────────┴──────────┐
                                                           │                     │
                                                     Access Token          Refresh Token
                                                     (15 min TTL)         (30 days TTL)
                                                           │                     │
                                                     Authorization          POST /auth/refresh
                                                     header                 (one-time use)
```

### OTP
- 4-digit numeric code
- TTL: 5 minutes
- Max attempts: 5 per session
- Stored in Redis: `otp:{session_id}` → `{phone, code, attempts, created_at}`
- Rate limit: 3 SMS/hour, 10 SMS/day per phone number
- Dev mode: code always "1234" when `SMS_PROVIDER=mock`

### JWT Tokens
```go
// Access Token (short-lived)
Claims: {
    "sub":  "user_uuid",
    "role": "user",         // user | admin | superadmin
    "iat":  1234567890,
    "exp":  1234568790,     // +15 min
}
Algorithm: HS256
Secret: JWT_SECRET env variable (min 32 bytes)

// Refresh Token (long-lived)
Claims: {
    "sub":  "user_uuid",
    "jti":  "token_uuid",  // unique token ID for revocation
    "iat":  1234567890,
    "exp":  1237159890,     // +30 days
}
```

### Refresh Token Rotation
```
1. Client sends POST /auth/refresh { refresh_token }
2. Server validates token, checks jti in Redis
3. Server deletes old jti from Redis (one-time use)
4. Server issues new access_token + new refresh_token
5. If old refresh_token is reused → revoke ALL user tokens (compromised)
```

### PIN Code
- 4-digit numeric
- Stored: bcrypt hash on device (react-native-keychain)
- Max 3 attempts → require re-auth via OTP
- Optional: used for quick re-login on same device
- NOT transmitted to server (local authentication only)

---

## 2. Authorization (RBAC)

### Global Roles
| Role | Access |
|------|--------|
| user | Own data, public data, community member functions |
| admin | Community admin panel (scoped to their communities) |
| superadmin | Superadmin panel, all data, court management |

### Community Roles
| Role | Permissions |
|------|------------|
| owner | All admin + delete community, transfer ownership |
| admin | Manage members, events, posts, settings, rating |
| moderator | Manage posts, approve events, moderate chat |
| member | View, participate, create events (if allowed), chat |

### Middleware Chain
```go
// Every request goes through:
r.Use(middleware.Logger)
r.Use(middleware.Recovery)
r.Use(middleware.CORS)
r.Use(middleware.RateLimit)
r.Use(middleware.RequestID)

// Protected routes add:
r.Use(middleware.Auth)              // validates JWT, sets user_id in context

// Admin routes add:
r.Use(middleware.RequireRole("admin", "superadmin"))

// Community-scoped routes add:
r.Use(middleware.RequireCommunityRole("admin", "moderator"))
```

### Implementation
```go
func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractBearerToken(r)
        if token == "" {
            respondError(w, 401, "UNAUTHORIZED", "Missing token")
            return
        }
        claims, err := jwt.Validate(token)
        if err != nil {
            respondError(w, 401, "UNAUTHORIZED", "Invalid token")
            return
        }
        ctx := context.WithValue(r.Context(), userIDKey, claims.Sub)
        ctx = context.WithValue(ctx, userRoleKey, claims.Role)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 3. Rate Limiting

### Implementation
Redis-based sliding window counter.

```go
func RateLimit(key string, limit int, window time.Duration) middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            identifier := resolveKey(key, r) // user_id or IP
            redisKey := fmt.Sprintf("rate:%s:%s", key, identifier)
            
            count, err := redis.Incr(ctx, redisKey)
            if count == 1 {
                redis.Expire(ctx, redisKey, window)
            }
            if count > limit {
                w.Header().Set("Retry-After", "60")
                respondError(w, 429, "RATE_LIMITED", "Too many requests")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### Limits

| Endpoint / Action | Limit | Window | Key |
|-------------------|-------|--------|-----|
| SMS send | 3 | 1 hour | phone |
| SMS send | 10 | 24 hours | phone |
| OTP verify | 5 | per session | session_id |
| PIN verify | 3 | per session | device_id |
| API general | 100 | 1 minute | user_id |
| Search | 30 | 1 minute | user_id |
| Event creation | 10 | 1 hour | user_id |
| Post creation | 20 | 1 hour | user_id |
| Message send | 60 | 1 minute | user_id |
| File upload | 10 | 1 hour | user_id |
| Registration | 3 | 1 hour | IP |

---

## 4. Input Validation

### Backend Validation (mandatory)
```go
// NEVER trust client input. Validate everything in handler layer.
type CreateEventRequest struct {
    Title           string  `json:"title" validate:"required,min=3,max=100"`
    Description     string  `json:"description" validate:"max=2000"`
    EventType       string  `json:"event_type" validate:"required,oneof=find_partner organized_game tournament training"`
    LevelMin        float64 `json:"level_min" validate:"required,gte=1.0,lte=7.0"`
    LevelMax        float64 `json:"level_max" validate:"required,gte=1.0,lte=7.0,gtefield=LevelMin"`
    MaxParticipants int     `json:"max_participants" validate:"required,gte=2,lte=64"`
    StartTime       string  `json:"start_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}
```

### Validation Rules

| Field | Rules |
|-------|-------|
| Phone | `^\+7[0-9]{10}$` (Kazakhstan format) |
| Name | 2-50 chars, letters + spaces + hyphens only |
| Title | 3-100 chars |
| Description | 0-2000 chars |
| Level (NTRP) | 1.0-7.0, step 0.5 |
| Rating | 100-3000 |
| PIN | exactly 4 digits |
| OTP | exactly 4 digits |
| UUID | valid UUID v4 format |
| URL | valid URL format (for avatars) |
| Pagination page | ≥ 1 |
| Pagination per_page | 1-100, default 20 |

---

## 5. Data Protection

### SQL Injection
- sqlc generates parameterized queries (safe by design)
- NEVER concatenate user input into SQL strings
- NEVER use `fmt.Sprintf` for SQL queries

### XSS
- Sanitize HTML in post content (bluemonday library)
- JSON responses auto-escaped by encoding/json
- Mobile: React Native auto-escapes text

### CORS
```go
cors.Options{
    AllowedOrigins:   []string{"https://admin.tennisapp.kz", "https://superadmin.tennisapp.kz"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
    AllowedHeaders:   []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
    MaxAge:           86400,
}
// Development: AllowedOrigins = ["http://localhost:*"]
```

### Sensitive Data
| Data | Protection |
|------|-----------|
| Phone numbers | Masked in logs: +7707***4567 |
| JWT secret | Environment variable only |
| PIN hash | Device keychain only (never sent to server) |
| OTP code | Redis only (5 min TTL), never in logs |
| Passwords | N/A (no passwords in system) |
| DB credentials | Environment variable only |
| API keys | Environment variable only |
| FCM tokens | Stored in DB, not exposed in API responses |

### HTTPS
- TLS 1.2+ required (enforced by Railway/Vercel)
- HSTS header enabled
- No mixed content

---

## 6. WebSocket Security

```go
// 1. JWT required for WebSocket connection
ws://api.tennisapp.kz/ws?token={jwt}

// 2. Validate token on upgrade
func (h *WSHandler) Upgrade(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    claims, err := jwt.Validate(token)
    if err != nil {
        http.Error(w, "Unauthorized", 401)
        return
    }
    // 3. Check user can access chat
    // 4. Rate limit messages (60/min)
}
```

---

## 7. File Upload Security

```go
func validateUpload(file multipart.File, header *multipart.FileHeader) error {
    // 1. Check file size (max 5 MB)
    if header.Size > 5*1024*1024 {
        return ErrFileTooLarge
    }
    // 2. Check MIME type by reading magic bytes (not just extension)
    buf := make([]byte, 512)
    file.Read(buf)
    file.Seek(0, 0)
    contentType := http.DetectContentType(buf)
    if !isAllowedType(contentType) { // image/jpeg, image/png, image/webp
        return ErrInvalidFileType
    }
    // 3. Generate unique filename (prevent path traversal)
    key := fmt.Sprintf("%s/%s%s", folder, uuid.New(), filepath.Ext(header.Filename))
    return nil
}
```

---

## 8. Security Checklist (per Sprint)

```
□ All new endpoints require authentication (middleware.Auth)
□ All user input validated (struct tags + service logic)
□ No SQL injection vectors (sqlc parameterized queries)
□ Rate limiting on new endpoints
□ RBAC enforced (community roles checked)
□ No secrets in code or logs
□ Error messages don't leak internal details
□ File uploads validated (type + size)
□ WebSocket messages authenticated
□ Pagination enforced (no unbounded queries)
```
