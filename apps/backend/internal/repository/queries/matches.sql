-- name: CreateMatch :one
INSERT INTO matches (
    event_id, community_id,
    player1_id, player2_id,
    player1_partner_id, player2_partner_id,
    composition,
    round_name, round_number, court_number,
    scheduled_time
) VALUES (
    sqlc.narg('event_id'), sqlc.narg('community_id'),
    @player1_id, @player2_id,
    sqlc.narg('player1_partner_id'), sqlc.narg('player2_partner_id'),
    @composition,
    sqlc.narg('round_name'), sqlc.narg('round_number'), sqlc.narg('court_number'),
    sqlc.narg('scheduled_time')
)
RETURNING id, event_id, community_id,
    player1_id, player2_id, player1_partner_id, player2_partner_id,
    composition, score, winner_id,
    result_status, submitted_by, confirmed_by, submitted_at, confirmed_at,
    dispute_reason,
    player1_rating_before, player1_rating_after,
    player2_rating_before, player2_rating_after,
    round_name, round_number, court_number,
    scheduled_time, played_at, created_at, updated_at;

-- name: GetMatchByID :one
SELECT id, event_id, community_id,
    player1_id, player2_id, player1_partner_id, player2_partner_id,
    composition, score, winner_id,
    result_status, submitted_by, confirmed_by, submitted_at, confirmed_at,
    dispute_reason,
    player1_rating_before, player1_rating_after,
    player2_rating_before, player2_rating_after,
    round_name, round_number, court_number,
    scheduled_time, played_at, created_at, updated_at
FROM matches
WHERE id = $1;

-- name: SubmitMatchResult :one
UPDATE matches SET
    score = @score,
    winner_id = @winner_id,
    submitted_by = @submitted_by,
    submitted_at = NOW(),
    result_status = 'pending',
    played_at = COALESCE(played_at, NOW()),
    updated_at = NOW()
WHERE id = @id
RETURNING id, event_id, community_id,
    player1_id, player2_id, player1_partner_id, player2_partner_id,
    composition, score, winner_id,
    result_status, submitted_by, confirmed_by, submitted_at, confirmed_at,
    dispute_reason,
    player1_rating_before, player1_rating_after,
    player2_rating_before, player2_rating_after,
    round_name, round_number, court_number,
    scheduled_time, played_at, created_at, updated_at;

-- name: ConfirmMatch :one
UPDATE matches SET
    result_status = 'confirmed',
    confirmed_by = @confirmed_by,
    confirmed_at = NOW(),
    player1_rating_before = @player1_rating_before,
    player1_rating_after = @player1_rating_after,
    player2_rating_before = @player2_rating_before,
    player2_rating_after = @player2_rating_after,
    updated_at = NOW()
WHERE id = @id
RETURNING id, event_id, community_id,
    player1_id, player2_id, player1_partner_id, player2_partner_id,
    composition, score, winner_id,
    result_status, submitted_by, confirmed_by, submitted_at, confirmed_at,
    dispute_reason,
    player1_rating_before, player1_rating_after,
    player2_rating_before, player2_rating_after,
    round_name, round_number, court_number,
    scheduled_time, played_at, created_at, updated_at;

-- name: DisputeMatch :one
UPDATE matches SET
    result_status = 'disputed',
    dispute_reason = @dispute_reason,
    updated_at = NOW()
WHERE id = @id
RETURNING id, event_id, community_id,
    player1_id, player2_id, player1_partner_id, player2_partner_id,
    composition, score, winner_id,
    result_status, submitted_by, confirmed_by, submitted_at, confirmed_at,
    dispute_reason,
    player1_rating_before, player1_rating_after,
    player2_rating_before, player2_rating_after,
    round_name, round_number, court_number,
    scheduled_time, played_at, created_at, updated_at;

-- name: AdminConfirmMatch :one
UPDATE matches SET
    result_status = 'admin_confirmed',
    score = @score,
    winner_id = @winner_id,
    confirmed_by = @confirmed_by,
    confirmed_at = NOW(),
    player1_rating_before = @player1_rating_before,
    player1_rating_after = @player1_rating_after,
    player2_rating_before = @player2_rating_before,
    player2_rating_after = @player2_rating_after,
    updated_at = NOW()
WHERE id = @id
RETURNING id, event_id, community_id,
    player1_id, player2_id, player1_partner_id, player2_partner_id,
    composition, score, winner_id,
    result_status, submitted_by, confirmed_by, submitted_at, confirmed_at,
    dispute_reason,
    player1_rating_before, player1_rating_after,
    player2_rating_before, player2_rating_after,
    round_name, round_number, court_number,
    scheduled_time, played_at, created_at, updated_at;

-- name: ListMyMatches :many
SELECT id, event_id, community_id,
    player1_id, player2_id, player1_partner_id, player2_partner_id,
    composition, score, winner_id,
    result_status, submitted_by, confirmed_by, submitted_at, confirmed_at,
    dispute_reason,
    player1_rating_before, player1_rating_after,
    player2_rating_before, player2_rating_after,
    round_name, round_number, court_number,
    scheduled_time, played_at, created_at, updated_at
FROM matches
WHERE (player1_id = @user_id OR player2_id = @user_id
       OR player1_partner_id = @user_id OR player2_partner_id = @user_id)
  AND (sqlc.narg('community_id')::uuid IS NULL OR community_id = sqlc.narg('community_id'))
  AND (sqlc.narg('opponent_id')::uuid IS NULL OR
       (player1_id = sqlc.narg('opponent_id') OR player2_id = sqlc.narg('opponent_id')))
  AND (sqlc.narg('result_filter')::text IS NULL
       OR (sqlc.narg('result_filter')::text = 'win' AND winner_id = @user_id)
       OR (sqlc.narg('result_filter')::text = 'loss' AND winner_id IS NOT NULL AND winner_id != @user_id))
ORDER BY COALESCE(played_at, created_at) DESC
LIMIT @result_limit OFFSET @result_offset;

-- name: CountMyMatches :one
SELECT COUNT(*)
FROM matches
WHERE (player1_id = @user_id OR player2_id = @user_id
       OR player1_partner_id = @user_id OR player2_partner_id = @user_id)
  AND (sqlc.narg('community_id')::uuid IS NULL OR community_id = sqlc.narg('community_id'))
  AND (sqlc.narg('opponent_id')::uuid IS NULL OR
       (player1_id = sqlc.narg('opponent_id') OR player2_id = sqlc.narg('opponent_id')))
  AND (sqlc.narg('result_filter')::text IS NULL
       OR (sqlc.narg('result_filter')::text = 'win' AND winner_id = @user_id)
       OR (sqlc.narg('result_filter')::text = 'loss' AND winner_id IS NOT NULL AND winner_id != @user_id));

-- name: UpdateUserRating :exec
UPDATE users SET
    global_rating = @new_rating,
    updated_at = NOW()
WHERE id = @user_id;

-- name: InsertRatingHistory :one
INSERT INTO rating_history (
    user_id, community_id, rating_before, rating_after, change, match_id, reason
) VALUES (
    @user_id, sqlc.narg('community_id'), @rating_before, @rating_after, @change, @match_id, @reason
)
RETURNING id, user_id, community_id, rating_before, rating_after, change, match_id, reason, created_at;

-- name: UpsertPlayerStatsGlobal :exec
INSERT INTO player_stats_global (
    user_id, total_games, total_wins, total_losses, win_rate,
    singles_games, singles_wins, doubles_games, doubles_wins,
    current_streak, best_streak, last_game_at, updated_at
) VALUES (
    @user_id, 1,
    CASE WHEN @is_winner::boolean THEN 1 ELSE 0 END,
    CASE WHEN @is_winner::boolean THEN 0 ELSE 1 END,
    CASE WHEN @is_winner::boolean THEN 100.00 ELSE 0.00 END,
    CASE WHEN @is_singles::boolean THEN 1 ELSE 0 END,
    CASE WHEN @is_singles::boolean AND @is_winner::boolean THEN 1 ELSE 0 END,
    CASE WHEN @is_singles::boolean THEN 0 ELSE 1 END,
    CASE WHEN NOT @is_singles::boolean AND @is_winner::boolean THEN 1 ELSE 0 END,
    CASE WHEN @is_winner::boolean THEN 1 ELSE 0 END,
    CASE WHEN @is_winner::boolean THEN 1 ELSE 0 END,
    NOW(), NOW()
)
ON CONFLICT (user_id) DO UPDATE SET
    total_games = player_stats_global.total_games + 1,
    total_wins = player_stats_global.total_wins + CASE WHEN @is_winner::boolean THEN 1 ELSE 0 END,
    total_losses = player_stats_global.total_losses + CASE WHEN @is_winner::boolean THEN 0 ELSE 1 END,
    win_rate = ROUND(
        (player_stats_global.total_wins + CASE WHEN @is_winner::boolean THEN 1 ELSE 0 END)::decimal /
        (player_stats_global.total_games + 1) * 100, 2
    ),
    singles_games = player_stats_global.singles_games + CASE WHEN @is_singles::boolean THEN 1 ELSE 0 END,
    singles_wins = player_stats_global.singles_wins + CASE WHEN @is_singles::boolean AND @is_winner::boolean THEN 1 ELSE 0 END,
    doubles_games = player_stats_global.doubles_games + CASE WHEN @is_singles::boolean THEN 0 ELSE 1 END,
    doubles_wins = player_stats_global.doubles_wins + CASE WHEN NOT @is_singles::boolean AND @is_winner::boolean THEN 1 ELSE 0 END,
    current_streak = CASE
        WHEN @is_winner::boolean THEN player_stats_global.current_streak + 1
        ELSE 0
    END,
    best_streak = GREATEST(
        player_stats_global.best_streak,
        CASE WHEN @is_winner::boolean THEN player_stats_global.current_streak + 1 ELSE 0 END
    ),
    last_game_at = NOW(),
    updated_at = NOW();

-- name: UpdateCommunityMemberStats :exec
UPDATE community_members SET
    community_rating = @new_rating,
    community_games_count = community_games_count + 1,
    community_wins = community_wins + CASE WHEN @is_winner::boolean THEN 1 ELSE 0 END,
    community_losses = community_losses + CASE WHEN @is_winner::boolean THEN 0 ELSE 1 END,
    updated_at = NOW()
WHERE community_id = @community_id AND user_id = @user_id;

-- name: GetGlobalLeaderboard :many
SELECT
    u.id as user_id,
    u.first_name, u.last_name, u.avatar_url, u.ntrp_level,
    u.global_rating,
    COALESCE(ps.total_games, 0) as total_games,
    COALESCE(ps.total_wins, 0) as total_wins,
    COALESCE(ps.total_losses, 0) as total_losses,
    COALESCE(ps.win_rate, 0) as win_rate
FROM users u
LEFT JOIN player_stats_global ps ON u.id = ps.user_id
WHERE u.status = 'active'
  AND u.global_rating IS NOT NULL
  AND (sqlc.narg('min_games')::int IS NULL OR COALESCE(ps.total_games, 0) >= sqlc.narg('min_games')::int)
ORDER BY u.global_rating DESC
LIMIT @result_limit OFFSET @result_offset;

-- name: CountGlobalLeaderboard :one
SELECT COUNT(*)
FROM users u
LEFT JOIN player_stats_global ps ON u.id = ps.user_id
WHERE u.status = 'active'
  AND u.global_rating IS NOT NULL
  AND (sqlc.narg('min_games')::int IS NULL OR COALESCE(ps.total_games, 0) >= sqlc.narg('min_games')::int);

-- name: GetRatingHistory :many
SELECT id, user_id, community_id, rating_before, rating_after, change, match_id, reason, created_at
FROM rating_history
WHERE user_id = @user_id
  AND (sqlc.narg('community_id')::uuid IS NULL OR community_id = sqlc.narg('community_id'))
  AND (sqlc.narg('since')::timestamptz IS NULL OR created_at >= sqlc.narg('since'))
ORDER BY created_at DESC
LIMIT @result_limit OFFSET @result_offset;

-- name: GetUserRatingPosition :one
SELECT
    u.global_rating,
    (SELECT COUNT(*) FROM users u2 WHERE u2.status = 'active' AND u2.global_rating > u.global_rating) + 1 as rank,
    (SELECT COUNT(*) FROM users u3 WHERE u3.status = 'active' AND u3.global_rating IS NOT NULL) as total_players
FROM users u
WHERE u.id = @user_id;

-- name: GetUserCommunityRatings :many
SELECT
    cm.community_id,
    c.name as community_name,
    c.logo_url as community_logo,
    cm.community_rating,
    cm.community_games_count,
    cm.community_wins,
    cm.community_losses,
    (SELECT COUNT(*) FROM community_members cm2
     WHERE cm2.community_id = cm.community_id AND cm2.status = 'active'
       AND cm2.community_rating > cm.community_rating) + 1 as rank
FROM community_members cm
JOIN communities c ON cm.community_id = c.id
WHERE cm.user_id = @user_id AND cm.status = 'active' AND c.is_active = TRUE
ORDER BY cm.community_rating DESC;

-- name: GetUserForRating :one
SELECT id, global_rating, ntrp_level
FROM users
WHERE id = $1;

-- name: UpdateUserNTRPLevel :exec
UPDATE users SET
    ntrp_level = @ntrp_level,
    updated_at = NOW()
WHERE id = @user_id;

-- name: GetPlayerTotalGames :one
SELECT COALESCE(total_games, 0)::int as total_games
FROM player_stats_global
WHERE user_id = $1;
