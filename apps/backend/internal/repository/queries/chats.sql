-- name: CreatePersonalChat :one
INSERT INTO chats (chat_type, user1_id, user2_id)
VALUES ('personal', @user1_id, @user2_id)
RETURNING id, chat_type, community_id, event_id, user1_id, user2_id,
    name, is_archived, last_message_at, last_message_preview,
    pinned_message_id, created_at, updated_at;

-- name: GetPersonalChat :one
SELECT id, chat_type, community_id, event_id, user1_id, user2_id,
    name, is_archived, last_message_at, last_message_preview,
    pinned_message_id, created_at, updated_at
FROM chats
WHERE chat_type = 'personal'
  AND user1_id = @user1_id AND user2_id = @user2_id;

-- name: CreateCommunityChat :one
INSERT INTO chats (chat_type, community_id, name)
VALUES ('community', @community_id, @name)
RETURNING id, chat_type, community_id, event_id, user1_id, user2_id,
    name, is_archived, last_message_at, last_message_preview,
    pinned_message_id, created_at, updated_at;

-- name: GetCommunityChatByCommunityID :one
SELECT id, chat_type, community_id, event_id, user1_id, user2_id,
    name, is_archived, last_message_at, last_message_preview,
    pinned_message_id, created_at, updated_at
FROM chats
WHERE chat_type = 'community' AND community_id = @community_id;

-- name: CreateEventChat :one
INSERT INTO chats (chat_type, event_id, name)
VALUES ('event', @event_id, @name)
RETURNING id, chat_type, community_id, event_id, user1_id, user2_id,
    name, is_archived, last_message_at, last_message_preview,
    pinned_message_id, created_at, updated_at;

-- name: GetEventChatByEventID :one
SELECT id, chat_type, community_id, event_id, user1_id, user2_id,
    name, is_archived, last_message_at, last_message_preview,
    pinned_message_id, created_at, updated_at
FROM chats
WHERE chat_type = 'event' AND event_id = @event_id;

-- name: GetChatByID :one
SELECT id, chat_type, community_id, event_id, user1_id, user2_id,
    name, is_archived, last_message_at, last_message_preview,
    pinned_message_id, created_at, updated_at
FROM chats
WHERE id = @id;

-- name: ListMyChats :many
SELECT
    c.id, c.chat_type, c.community_id, c.event_id, c.user1_id, c.user2_id,
    c.name, c.is_archived, c.last_message_at, c.last_message_preview,
    c.pinned_message_id, c.created_at, c.updated_at,
    COALESCE(crs.is_muted, FALSE) as is_muted,
    crs.last_read_at,
    -- Count unread messages
    (SELECT COUNT(*) FROM messages m
     WHERE m.chat_id = c.id
       AND m.is_deleted = FALSE
       AND m.sender_id != @user_id
       AND (crs.last_read_at IS NULL OR m.created_at > crs.last_read_at)
    )::int as unread_count,
    -- Other user info for personal chats
    CASE WHEN c.chat_type = 'personal' AND c.user1_id = @user_id THEN c.user2_id
         WHEN c.chat_type = 'personal' AND c.user2_id = @user_id THEN c.user1_id
         ELSE NULL END as other_user_id
FROM chats c
LEFT JOIN chat_read_status crs ON c.id = crs.chat_id AND crs.user_id = @user_id
WHERE (
    -- Personal chats: user is user1 or user2
    (c.chat_type = 'personal' AND (c.user1_id = @user_id OR c.user2_id = @user_id))
    -- Community chats: user is member of community
    OR (c.chat_type = 'community' AND EXISTS (
        SELECT 1 FROM community_members cm
        WHERE cm.community_id = c.community_id AND cm.user_id = @user_id AND cm.status = 'active'
    ))
    -- Event chats: user is participant
    OR (c.chat_type = 'event' AND EXISTS (
        SELECT 1 FROM event_participants ep
        WHERE ep.event_id = c.event_id AND ep.user_id = @user_id AND ep.status != 'cancelled'
    ))
)
AND c.is_archived = FALSE
ORDER BY COALESCE(c.last_message_at, c.created_at) DESC;

-- name: CreateMessage :one
INSERT INTO messages (chat_id, sender_id, content, reply_to_id)
VALUES (@chat_id, @sender_id, @content, sqlc.narg('reply_to_id'))
RETURNING id, chat_id, sender_id, content, reply_to_id, is_deleted, deleted_at, created_at;

-- name: GetMessageByID :one
SELECT id, chat_id, sender_id, content, reply_to_id, is_deleted, deleted_at, created_at
FROM messages
WHERE id = @id AND is_deleted = FALSE;

-- name: GetMessages :many
SELECT
    m.id, m.chat_id, m.sender_id, m.content, m.reply_to_id,
    m.is_deleted, m.deleted_at, m.created_at,
    u.first_name as sender_first_name,
    u.last_name as sender_last_name,
    u.avatar_url as sender_avatar_url
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.chat_id = @chat_id
  AND m.is_deleted = FALSE
  AND (sqlc.narg('before_id')::uuid IS NULL OR m.id < sqlc.narg('before_id'))
ORDER BY m.created_at DESC
LIMIT @result_limit;

-- name: UpdateChatLastMessage :exec
UPDATE chats SET
    last_message_at = @last_message_at,
    last_message_preview = @preview,
    updated_at = NOW()
WHERE id = @id;

-- name: UpsertChatReadStatus :exec
INSERT INTO chat_read_status (chat_id, user_id, last_read_at)
VALUES (@chat_id, @user_id, @last_read_at)
ON CONFLICT (chat_id, user_id) DO UPDATE SET
    last_read_at = @last_read_at;

-- name: UpdateChatMuted :exec
INSERT INTO chat_read_status (chat_id, user_id, is_muted)
VALUES (@chat_id, @user_id, @is_muted)
ON CONFLICT (chat_id, user_id) DO UPDATE SET
    is_muted = @is_muted;

-- name: GetTotalUnreadCount :one
SELECT COALESCE(SUM(unread), 0)::int as total_unread
FROM (
    SELECT
        (SELECT COUNT(*) FROM messages m
         WHERE m.chat_id = c.id
           AND m.is_deleted = FALSE
           AND m.sender_id != @user_id
           AND (crs.last_read_at IS NULL OR m.created_at > crs.last_read_at)
        ) as unread
    FROM chats c
    LEFT JOIN chat_read_status crs ON c.id = crs.chat_id AND crs.user_id = @user_id
    WHERE (
        (c.chat_type = 'personal' AND (c.user1_id = @user_id OR c.user2_id = @user_id))
        OR (c.chat_type = 'community' AND EXISTS (
            SELECT 1 FROM community_members cm
            WHERE cm.community_id = c.community_id AND cm.user_id = @user_id AND cm.status = 'active'
        ))
        OR (c.chat_type = 'event' AND EXISTS (
            SELECT 1 FROM event_participants ep
            WHERE ep.event_id = c.event_id AND ep.user_id = @user_id AND ep.status != 'cancelled'
        ))
    )
    AND c.is_archived = FALSE
    AND COALESCE(crs.is_muted, FALSE) = FALSE
) sub;

-- name: IsUserInChat :one
SELECT EXISTS (
    SELECT 1 FROM chats c
    WHERE c.id = @chat_id
    AND (
        (c.chat_type = 'personal' AND (c.user1_id = @user_id OR c.user2_id = @user_id))
        OR (c.chat_type = 'community' AND EXISTS (
            SELECT 1 FROM community_members cm
            WHERE cm.community_id = c.community_id AND cm.user_id = @user_id AND cm.status = 'active'
        ))
        OR (c.chat_type = 'event' AND EXISTS (
            SELECT 1 FROM event_participants ep
            WHERE ep.event_id = c.event_id AND ep.user_id = @user_id AND ep.status != 'cancelled'
        ))
    )
) as is_member;

-- name: GetUserBasicInfo :one
SELECT id, first_name, last_name, avatar_url
FROM users
WHERE id = @id AND status = 'active';

-- name: GetCommunityBasicInfo :one
SELECT id, name, logo_url
FROM communities
WHERE id = @id AND is_active = TRUE;

-- name: GetEventBasicInfo :one
SELECT id, title
FROM events
WHERE id = @id;

-- name: GetChatMembersForPersonal :many
SELECT user1_id, user2_id
FROM chats
WHERE id = @id AND chat_type = 'personal';

-- name: GetChatMembersForCommunity :many
SELECT cm.user_id
FROM community_members cm
JOIN chats c ON c.community_id = cm.community_id
WHERE c.id = @chat_id AND cm.status = 'active';

-- name: GetChatMembersForEvent :many
SELECT ep.user_id
FROM event_participants ep
JOIN chats c ON c.event_id = ep.event_id
WHERE c.id = @chat_id AND ep.status != 'cancelled';
