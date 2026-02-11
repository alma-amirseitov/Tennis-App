-- name: CreateEvent :one
INSERT INTO events (
    title, description, event_type, status,
    community_id, player_composition, match_format, match_format_details,
    court_id, location_name, location_address,
    start_time, end_time,
    max_participants, min_participants,
    min_level, max_level,
    gender_restriction, registration_deadline,
    is_paid, price_amount, price_currency,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
    $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
)
RETURNING id, title, description, event_type, status,
    community_id, player_composition, match_format, match_format_details,
    tournament_system, tournament_details,
    court_id, location_name, location_address,
    start_time, end_time,
    max_participants, min_participants, current_participants,
    min_level, max_level,
    gender_restriction, min_age, max_age,
    registration_deadline,
    is_paid, price_amount, price_currency,
    created_by, created_at, updated_at;

-- name: GetEventByID :one
SELECT id, title, description, event_type, status,
    community_id, player_composition, match_format, match_format_details,
    tournament_system, tournament_details,
    court_id, location_name, location_address,
    start_time, end_time,
    max_participants, min_participants, current_participants,
    min_level, max_level,
    gender_restriction, min_age, max_age,
    registration_deadline,
    is_paid, price_amount, price_currency,
    created_by, created_at, updated_at
FROM events
WHERE id = $1;

-- name: ListEvents :many
SELECT id, title, description, event_type, status,
    community_id, player_composition, match_format,
    court_id, location_name, location_address,
    start_time, end_time,
    max_participants, min_participants, current_participants,
    min_level, max_level,
    gender_restriction, registration_deadline,
    is_paid, price_amount, price_currency,
    created_by, created_at
FROM events
WHERE (sqlc.narg('event_type')::event_type IS NULL OR event_type = sqlc.narg('event_type'))
  AND (sqlc.narg('event_status')::event_status IS NULL OR status = sqlc.narg('event_status'))
  AND (sqlc.narg('composition')::player_composition IS NULL OR player_composition = sqlc.narg('composition'))
  AND (sqlc.narg('community_id')::uuid IS NULL OR community_id = sqlc.narg('community_id'))
  AND (sqlc.narg('min_level')::decimal IS NULL OR max_level >= sqlc.narg('min_level'))
  AND (sqlc.narg('max_level')::decimal IS NULL OR min_level <= sqlc.narg('max_level'))
  AND (sqlc.narg('date_from')::timestamptz IS NULL OR start_time >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR start_time <= sqlc.narg('date_to'))
  AND (sqlc.narg('district')::text IS NULL OR location_address ILIKE '%' || sqlc.narg('district') || '%')
  AND status NOT IN ('draft', 'cancelled', 'archived')
ORDER BY
    CASE WHEN @sort_by::text = 'date_asc' THEN start_time END ASC,
    CASE WHEN @sort_by::text = 'spots_left' THEN (max_participants - current_participants) END DESC,
    start_time DESC
LIMIT @result_limit OFFSET @result_offset;

-- name: CountEvents :one
SELECT COUNT(*)
FROM events
WHERE (sqlc.narg('event_type')::event_type IS NULL OR event_type = sqlc.narg('event_type'))
  AND (sqlc.narg('event_status')::event_status IS NULL OR status = sqlc.narg('event_status'))
  AND (sqlc.narg('composition')::player_composition IS NULL OR player_composition = sqlc.narg('composition'))
  AND (sqlc.narg('community_id')::uuid IS NULL OR community_id = sqlc.narg('community_id'))
  AND (sqlc.narg('min_level')::decimal IS NULL OR max_level >= sqlc.narg('min_level'))
  AND (sqlc.narg('max_level')::decimal IS NULL OR min_level <= sqlc.narg('max_level'))
  AND (sqlc.narg('date_from')::timestamptz IS NULL OR start_time >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR start_time <= sqlc.narg('date_to'))
  AND (sqlc.narg('district')::text IS NULL OR location_address ILIKE '%' || sqlc.narg('district') || '%')
  AND status NOT IN ('draft', 'cancelled', 'archived');

-- name: UpdateEventStatus :one
UPDATE events SET
    status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, title, description, event_type, status,
    community_id, player_composition, match_format, match_format_details,
    court_id, location_name, location_address,
    start_time, end_time,
    max_participants, min_participants, current_participants,
    min_level, max_level,
    gender_restriction, registration_deadline,
    is_paid, price_amount, price_currency,
    created_by, created_at, updated_at;

-- name: UpdateEvent :one
UPDATE events SET
    title              = COALESCE(sqlc.narg('title'), title),
    description        = COALESCE(sqlc.narg('description'), description),
    start_time         = COALESCE(sqlc.narg('start_time'), start_time),
    end_time           = COALESCE(sqlc.narg('end_time'), end_time),
    max_participants   = COALESCE(sqlc.narg('max_participants'), max_participants),
    min_level          = COALESCE(sqlc.narg('min_level'), min_level),
    max_level          = COALESCE(sqlc.narg('max_level'), max_level),
    registration_deadline = COALESCE(sqlc.narg('registration_deadline'), registration_deadline),
    updated_at         = NOW()
WHERE id = @id
RETURNING id, title, description, event_type, status,
    community_id, player_composition, match_format, match_format_details,
    court_id, location_name, location_address,
    start_time, end_time,
    max_participants, min_participants, current_participants,
    min_level, max_level,
    gender_restriction, registration_deadline,
    is_paid, price_amount, price_currency,
    created_by, created_at, updated_at;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;

-- name: AddEventParticipant :one
INSERT INTO event_participants (
    event_id, user_id, status, partner_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, event_id, user_id, status, registered_at, cancelled_at, partner_id, seed_number;

-- name: GetEventParticipant :one
SELECT id, event_id, user_id, status, registered_at, cancelled_at, partner_id, seed_number
FROM event_participants
WHERE event_id = $1 AND user_id = $2;

-- name: RemoveEventParticipant :exec
DELETE FROM event_participants
WHERE event_id = $1 AND user_id = $2;

-- name: ListEventParticipants :many
SELECT ep.id, ep.event_id, ep.user_id, ep.status, ep.registered_at, ep.partner_id, ep.seed_number,
    u.first_name, u.last_name, u.avatar_url, u.ntrp_level, u.global_rating
FROM event_participants ep
JOIN users u ON ep.user_id = u.id
WHERE ep.event_id = $1 AND ep.status IN ('registered', 'confirmed', 'checked_in')
ORDER BY ep.registered_at ASC;

-- name: GetCalendarEvents :many
SELECT id, title, event_type, status, start_time, end_time,
    community_id, location_name, current_participants, max_participants
FROM events
WHERE start_time >= @month_start::timestamptz
  AND start_time < @month_end::timestamptz
  AND (sqlc.narg('community_id')::uuid IS NULL OR community_id = sqlc.narg('community_id'))
  AND status NOT IN ('draft', 'cancelled', 'archived')
ORDER BY start_time ASC;

-- name: ListMyCreatedEvents :many
SELECT id, title, description, event_type, status,
    community_id, start_time, end_time,
    max_participants, current_participants,
    min_level, max_level, location_name,
    created_by, created_at
FROM events
WHERE created_by = $1
ORDER BY start_time DESC
LIMIT $2 OFFSET $3;

-- name: ListMyJoinedEvents :many
SELECT e.id, e.title, e.description, e.event_type, e.status,
    e.community_id, e.start_time, e.end_time,
    e.max_participants, e.current_participants,
    e.min_level, e.max_level, e.location_name,
    e.created_by, e.created_at
FROM events e
JOIN event_participants ep ON e.id = ep.event_id
WHERE ep.user_id = $1
  AND ep.status IN ('registered', 'confirmed', 'checked_in')
  AND e.status NOT IN ('completed', 'cancelled', 'archived')
ORDER BY e.start_time ASC
LIMIT $2 OFFSET $3;

-- name: ListMyPastEvents :many
SELECT e.id, e.title, e.description, e.event_type, e.status,
    e.community_id, e.start_time, e.end_time,
    e.max_participants, e.current_participants,
    e.min_level, e.max_level, e.location_name,
    e.created_by, e.created_at
FROM events e
JOIN event_participants ep ON e.id = ep.event_id
WHERE ep.user_id = $1
  AND e.status IN ('completed', 'cancelled', 'archived')
ORDER BY e.start_time DESC
LIMIT $2 OFFSET $3;

-- name: GetUserStats :one
SELECT user_id, total_games, total_wins, total_losses, win_rate,
    singles_games, singles_wins, doubles_games, doubles_wins,
    current_streak, best_streak, tournaments_played,
    last_game_at, updated_at
FROM player_stats_global
WHERE user_id = $1;

-- name: GetUserBadges :many
SELECT ub.badge_id, ub.earned_at,
    bd.name_ru, bd.name_kz, bd.name_en,
    bd.description_ru, bd.icon, bd.condition_type, bd.condition_value
FROM user_badges ub
JOIN badge_definitions bd ON ub.badge_id = bd.id
WHERE ub.user_id = $1
ORDER BY ub.earned_at DESC;

-- name: GetUserCommunities :many
SELECT c.id, c.name, c.slug, c.logo_url, c.community_type,
    cm.role, cm.status
FROM community_members cm
JOIN communities c ON cm.community_id = c.id
WHERE cm.user_id = $1 AND cm.status = 'active' AND c.is_active = TRUE
ORDER BY cm.joined_at DESC;

-- name: CheckFriendship :one
SELECT EXISTS(
    SELECT 1 FROM friends WHERE user_id = $1 AND friend_id = $2
) as is_friend;

-- name: CountMutualCommunities :one
SELECT COUNT(DISTINCT cm1.community_id)
FROM community_members cm1
JOIN community_members cm2 ON cm1.community_id = cm2.community_id
WHERE cm1.user_id = $1 AND cm2.user_id = $2
  AND cm1.status = 'active' AND cm2.status = 'active';

-- name: UpdateUserAvatarURL :one
UPDATE users SET
    avatar_url = $2,
    updated_at = NOW()
WHERE id = $1 AND status != 'deleted'
RETURNING id, avatar_url;
