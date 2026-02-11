-- name: CreateUser :one
INSERT INTO users (
    phone
) VALUES (
    $1
)
RETURNING id, phone, phone_verified,
    first_name, last_name, gender, birth_year, city, district, avatar_url, bio,
    ntrp_level, level_label, quiz_completed,
    global_rating, global_games_count,
    language, pin_hash, push_token, platform_role,
    profile_visibility, allow_messages_from, show_stats,
    notification_settings,
    status, is_profile_complete, last_active_at,
    created_at, updated_at;

-- name: GetUserByID :one
SELECT id, phone, phone_verified,
    first_name, last_name, gender, birth_year, city, district, avatar_url, bio,
    ntrp_level, level_label, quiz_completed,
    global_rating, global_games_count,
    language, pin_hash, push_token, platform_role,
    profile_visibility, allow_messages_from, show_stats,
    notification_settings,
    status, is_profile_complete, last_active_at,
    created_at, updated_at
FROM users
WHERE id = $1 AND status != 'deleted';

-- name: GetUserByPhone :one
SELECT id, phone, phone_verified,
    first_name, last_name, gender, birth_year, city, district, avatar_url, bio,
    ntrp_level, level_label, quiz_completed,
    global_rating, global_games_count,
    language, pin_hash, push_token, platform_role,
    profile_visibility, allow_messages_from, show_stats,
    notification_settings,
    status, is_profile_complete, last_active_at,
    created_at, updated_at
FROM users
WHERE phone = $1 AND status != 'deleted';

-- name: UpdateUser :one
UPDATE users SET
    first_name           = COALESCE(sqlc.narg('first_name'), first_name),
    last_name            = COALESCE(sqlc.narg('last_name'), last_name),
    gender               = COALESCE(sqlc.narg('gender'), gender),
    birth_year           = COALESCE(sqlc.narg('birth_year'), birth_year),
    city                 = COALESCE(sqlc.narg('city'), city),
    district             = COALESCE(sqlc.narg('district'), district),
    avatar_url           = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    bio                  = COALESCE(sqlc.narg('bio'), bio),
    ntrp_level           = COALESCE(sqlc.narg('ntrp_level'), ntrp_level),
    level_label          = COALESCE(sqlc.narg('level_label'), level_label),
    quiz_completed       = COALESCE(sqlc.narg('quiz_completed'), quiz_completed),
    global_rating        = COALESCE(sqlc.narg('global_rating'), global_rating),
    global_games_count   = COALESCE(sqlc.narg('global_games_count'), global_games_count),
    language             = COALESCE(sqlc.narg('language'), language),
    pin_hash             = COALESCE(sqlc.narg('pin_hash'), pin_hash),
    push_token           = COALESCE(sqlc.narg('push_token'), push_token),
    profile_visibility   = COALESCE(sqlc.narg('profile_visibility'), profile_visibility),
    allow_messages_from  = COALESCE(sqlc.narg('allow_messages_from'), allow_messages_from),
    show_stats           = COALESCE(sqlc.narg('show_stats'), show_stats),
    notification_settings = COALESCE(sqlc.narg('notification_settings'), notification_settings),
    is_profile_complete  = COALESCE(sqlc.narg('is_profile_complete'), is_profile_complete),
    last_active_at       = COALESCE(sqlc.narg('last_active_at'), last_active_at),
    updated_at           = NOW()
WHERE id = @id AND status != 'deleted'
RETURNING id, phone, phone_verified,
    first_name, last_name, gender, birth_year, city, district, avatar_url, bio,
    ntrp_level, level_label, quiz_completed,
    global_rating, global_games_count,
    language, pin_hash, push_token, platform_role,
    profile_visibility, allow_messages_from, show_stats,
    notification_settings,
    status, is_profile_complete, last_active_at,
    created_at, updated_at;

-- name: SearchUsers :many
SELECT id, phone, phone_verified,
    first_name, last_name, gender, birth_year, city, district, avatar_url, bio,
    ntrp_level, level_label, quiz_completed,
    global_rating, global_games_count,
    language, platform_role,
    profile_visibility, show_stats,
    status, is_profile_complete, last_active_at,
    created_at, updated_at
FROM users
WHERE status = 'active'
  AND is_profile_complete = TRUE
  AND (sqlc.narg('min_level')::decimal IS NULL OR ntrp_level >= sqlc.narg('min_level'))
  AND (sqlc.narg('max_level')::decimal IS NULL OR ntrp_level <= sqlc.narg('max_level'))
  AND (sqlc.narg('district')::text IS NULL OR district = sqlc.narg('district'))
  AND (sqlc.narg('gender')::gender_type IS NULL OR gender = sqlc.narg('gender'))
  AND (sqlc.narg('query')::text IS NULL OR (first_name || ' ' || last_name) ILIKE '%' || sqlc.narg('query') || '%')
ORDER BY
    CASE WHEN @sort_by::text = 'rating' THEN global_rating END DESC,
    CASE WHEN @sort_by::text = 'name' THEN first_name END ASC,
    CASE WHEN @sort_by::text = 'games' THEN global_games_count END DESC,
    CASE WHEN @sort_by::text = 'activity' THEN last_active_at END DESC NULLS LAST,
    global_rating DESC
LIMIT @result_limit OFFSET @result_offset;

-- name: CountSearchUsers :one
SELECT COUNT(*)
FROM users
WHERE status = 'active'
  AND is_profile_complete = TRUE
  AND (sqlc.narg('min_level')::decimal IS NULL OR ntrp_level >= sqlc.narg('min_level'))
  AND (sqlc.narg('max_level')::decimal IS NULL OR ntrp_level <= sqlc.narg('max_level'))
  AND (sqlc.narg('district')::text IS NULL OR district = sqlc.narg('district'))
  AND (sqlc.narg('gender')::gender_type IS NULL OR gender = sqlc.narg('gender'))
  AND (sqlc.narg('query')::text IS NULL OR (first_name || ' ' || last_name) ILIKE '%' || sqlc.narg('query') || '%');
