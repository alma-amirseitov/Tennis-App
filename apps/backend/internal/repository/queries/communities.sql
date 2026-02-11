-- name: CreateCommunity :one
INSERT INTO communities (
    name, slug, description, community_type, access_level,
    verification_status, district, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, name, slug, description, rules, community_type, access_level,
    verification_status, verified_at, verification_documents,
    logo_url, banner_url, contact_phone, contact_email, social_links,
    address, district,
    rating_initial, rating_k_factor, rating_min_games,
    member_count, event_count, is_active,
    created_by, created_at, updated_at;

-- name: GetCommunityByID :one
SELECT id, name, slug, description, rules, community_type, access_level,
    verification_status, verified_at, verification_documents,
    logo_url, banner_url, contact_phone, contact_email, social_links,
    address, district,
    rating_initial, rating_k_factor, rating_min_games,
    member_count, event_count, is_active,
    created_by, created_at, updated_at
FROM communities
WHERE id = $1 AND is_active = TRUE;

-- name: GetCommunityBySlug :one
SELECT id, name, slug, description, rules, community_type, access_level,
    verification_status, verified_at, verification_documents,
    logo_url, banner_url, contact_phone, contact_email, social_links,
    address, district,
    rating_initial, rating_k_factor, rating_min_games,
    member_count, event_count, is_active,
    created_by, created_at, updated_at
FROM communities
WHERE slug = $1 AND is_active = TRUE;

-- name: ListCommunities :many
SELECT id, name, slug, description, community_type, access_level,
    verification_status, logo_url, district,
    member_count, event_count, is_active,
    created_by, created_at
FROM communities
WHERE is_active = TRUE
  AND (sqlc.narg('community_type')::community_type IS NULL OR community_type = sqlc.narg('community_type'))
  AND (sqlc.narg('access_level')::community_access IS NULL OR access_level = sqlc.narg('access_level'))
  AND (sqlc.narg('verified_only')::boolean IS NULL OR (sqlc.narg('verified_only')::boolean = FALSE) OR verification_status = 'verified')
  AND (sqlc.narg('district')::text IS NULL OR district = sqlc.narg('district'))
  AND (sqlc.narg('query')::text IS NULL OR name ILIKE '%' || sqlc.narg('query') || '%')
ORDER BY
    CASE WHEN @sort_by::text = 'members' THEN member_count END DESC,
    CASE WHEN @sort_by::text = 'activity' THEN event_count END DESC,
    CASE WHEN @sort_by::text = 'name' THEN 1 END ASC,
    CASE WHEN @sort_by::text = 'created' THEN 1 END ASC,
    created_at DESC
LIMIT @result_limit OFFSET @result_offset;

-- name: CountCommunities :one
SELECT COUNT(*)
FROM communities
WHERE is_active = TRUE
  AND (sqlc.narg('community_type')::community_type IS NULL OR community_type = sqlc.narg('community_type'))
  AND (sqlc.narg('access_level')::community_access IS NULL OR access_level = sqlc.narg('access_level'))
  AND (sqlc.narg('verified_only')::boolean IS NULL OR (sqlc.narg('verified_only')::boolean = FALSE) OR verification_status = 'verified')
  AND (sqlc.narg('district')::text IS NULL OR district = sqlc.narg('district'))
  AND (sqlc.narg('query')::text IS NULL OR name ILIKE '%' || sqlc.narg('query') || '%');

-- name: UpdateCommunity :one
UPDATE communities SET
    name              = COALESCE(sqlc.narg('name'), name),
    description       = COALESCE(sqlc.narg('description'), description),
    rules             = COALESCE(sqlc.narg('rules'), rules),
    access_level      = COALESCE(sqlc.narg('access_level'), access_level),
    logo_url          = COALESCE(sqlc.narg('logo_url'), logo_url),
    banner_url        = COALESCE(sqlc.narg('banner_url'), banner_url),
    contact_phone     = COALESCE(sqlc.narg('contact_phone'), contact_phone),
    contact_email     = COALESCE(sqlc.narg('contact_email'), contact_email),
    social_links      = COALESCE(sqlc.narg('social_links'), social_links),
    address           = COALESCE(sqlc.narg('address'), address),
    district          = COALESCE(sqlc.narg('district'), district),
    updated_at        = NOW()
WHERE id = @id AND is_active = TRUE
RETURNING id, name, slug, description, rules, community_type, access_level,
    verification_status, verified_at, verification_documents,
    logo_url, banner_url, contact_phone, contact_email, social_links,
    address, district,
    rating_initial, rating_k_factor, rating_min_games,
    member_count, event_count, is_active,
    created_by, created_at, updated_at;

-- name: ListMyCommunities :many
SELECT c.id, c.name, c.slug, c.description, c.community_type, c.access_level,
    c.verification_status, c.logo_url, c.district,
    c.member_count, c.event_count, c.is_active,
    c.created_by, c.created_at,
    cm.role, cm.status as member_status
FROM communities c
JOIN community_members cm ON c.id = cm.community_id
WHERE cm.user_id = $1 AND cm.status = 'active' AND c.is_active = TRUE
ORDER BY cm.joined_at DESC;

-- name: AddCommunityMember :one
INSERT INTO community_members (
    community_id, user_id, role, status, application_message
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, community_id, user_id, role, status, application_message,
    reviewed_by, reviewed_at,
    community_rating, community_games_count, community_wins, community_losses,
    joined_at, updated_at;

-- name: GetCommunityMember :one
SELECT id, community_id, user_id, role, status, application_message,
    reviewed_by, reviewed_at,
    community_rating, community_games_count, community_wins, community_losses,
    joined_at, updated_at
FROM community_members
WHERE community_id = $1 AND user_id = $2;

-- name: UpdateCommunityMemberRole :one
UPDATE community_members SET
    role = $3,
    updated_at = NOW()
WHERE community_id = $1 AND user_id = $2
RETURNING id, community_id, user_id, role, status, application_message,
    reviewed_by, reviewed_at,
    community_rating, community_games_count, community_wins, community_losses,
    joined_at, updated_at;

-- name: UpdateCommunityMemberStatus :one
UPDATE community_members SET
    status = $3,
    reviewed_by = $4,
    reviewed_at = NOW(),
    updated_at = NOW()
WHERE community_id = $1 AND user_id = $2
RETURNING id, community_id, user_id, role, status, application_message,
    reviewed_by, reviewed_at,
    community_rating, community_games_count, community_wins, community_losses,
    joined_at, updated_at;

-- name: DeleteCommunityMember :exec
DELETE FROM community_members
WHERE community_id = $1 AND user_id = $2;

-- name: ListCommunityMembers :many
SELECT cm.id, cm.community_id, cm.user_id, cm.role, cm.status,
    cm.application_message, cm.reviewed_by, cm.reviewed_at,
    cm.community_rating, cm.community_games_count, cm.community_wins, cm.community_losses,
    cm.joined_at, cm.updated_at,
    u.first_name, u.last_name, u.avatar_url, u.ntrp_level, u.global_rating
FROM community_members cm
JOIN users u ON cm.user_id = u.id
WHERE cm.community_id = $1
  AND (sqlc.narg('member_role')::community_role IS NULL OR cm.role = sqlc.narg('member_role'))
  AND (sqlc.narg('member_status')::member_status IS NULL OR cm.status = sqlc.narg('member_status'))
  AND (sqlc.narg('query')::text IS NULL OR (u.first_name || ' ' || u.last_name) ILIKE '%' || sqlc.narg('query') || '%')
ORDER BY
    CASE WHEN @sort_by::text = 'rating' THEN cm.community_rating END DESC,
    CASE WHEN @sort_by::text = 'name' THEN u.first_name END ASC,
    CASE WHEN @sort_by::text = 'joined' THEN 1 END ASC,
    cm.joined_at DESC
LIMIT @result_limit OFFSET @result_offset;

-- name: CountCommunityMembers :one
SELECT COUNT(*)
FROM community_members cm
JOIN users u ON cm.user_id = u.id
WHERE cm.community_id = $1
  AND (sqlc.narg('member_role')::community_role IS NULL OR cm.role = sqlc.narg('member_role'))
  AND (sqlc.narg('member_status')::member_status IS NULL OR cm.status = sqlc.narg('member_status'))
  AND (sqlc.narg('query')::text IS NULL OR (u.first_name || ' ' || u.last_name) ILIKE '%' || sqlc.narg('query') || '%');

-- name: GetCommunityLeaderboard :many
SELECT
    cm.user_id,
    u.first_name, u.last_name, u.avatar_url, u.ntrp_level,
    cm.community_rating,
    cm.community_games_count,
    cm.community_wins,
    cm.community_losses,
    CASE WHEN cm.community_games_count > 0
        THEN ROUND(cm.community_wins::decimal / cm.community_games_count * 100, 1)
        ELSE 0
    END as win_rate
FROM community_members cm
JOIN users u ON cm.user_id = u.id
WHERE cm.community_id = $1 AND cm.status = 'active' AND u.status = 'active'
  AND cm.community_games_count >= $2
ORDER BY cm.community_rating DESC
LIMIT $3 OFFSET $4;
