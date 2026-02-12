-- Notifications queries

-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    type,
    title,
    body,
    data
) VALUES (
    @user_id,
    @type,
    @title,
    @body,
    sqlc.narg('data')
)
RETURNING
    id,
    user_id,
    type,
    title,
    body,
    data,
    is_read,
    is_pushed,
    created_at,
    expires_at;

-- name: ListNotifications :many
SELECT
    id,
    user_id,
    type,
    title,
    body,
    data,
    is_read,
    is_pushed,
    created_at,
    expires_at
FROM notifications
WHERE user_id = @user_id
ORDER BY created_at DESC
LIMIT @result_limit OFFSET @result_offset;

-- name: CountNotifications :one
SELECT COUNT(*)
FROM notifications
WHERE user_id = @user_id;

-- name: MarkNotificationRead :exec
UPDATE notifications
SET is_read = TRUE
WHERE id = @id AND user_id = @user_id;

-- name: MarkAllNotificationsRead :exec
UPDATE notifications
SET is_read = TRUE
WHERE user_id = @user_id AND is_read = FALSE;

-- name: GetUnreadNotificationCount :one
SELECT COUNT(*)
FROM notifications
WHERE user_id = @user_id AND is_read = FALSE;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = @id AND user_id = @user_id;

