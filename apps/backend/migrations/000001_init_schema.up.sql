-- =====================================================
-- Tennis Platform Database Schema v2.0
-- PostgreSQL 16+
-- Migration: 000001_init_schema
-- =====================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- =====================================================
-- ENUMS
-- =====================================================

CREATE TYPE gender_type AS ENUM ('male', 'female', 'other');
CREATE TYPE user_status AS ENUM ('active', 'banned', 'deleted');
CREATE TYPE platform_role AS ENUM ('player', 'superadmin');

CREATE TYPE community_type AS ENUM ('club', 'league', 'organizer', 'group');
CREATE TYPE community_access AS ENUM ('open', 'closed', 'paid');
CREATE TYPE verification_status AS ENUM ('none', 'pending', 'verified', 'rejected');
CREATE TYPE community_role AS ENUM ('owner', 'admin', 'moderator', 'coach_referee', 'member');
CREATE TYPE member_status AS ENUM ('pending', 'active', 'banned', 'left');

CREATE TYPE event_type AS ENUM ('find_partner', 'organized_game', 'tournament', 'training');
CREATE TYPE event_status AS ENUM (
    'draft', 'published', 'registration_open', 'registration_closed',
    'in_progress', 'completed', 'cancelled', 'archived'
);
CREATE TYPE player_composition AS ENUM ('singles', 'doubles', 'mixed', 'team', 'custom');
CREATE TYPE match_format AS ENUM ('best_of', 'pro_set', 'short_set', 'timed', 'custom');
CREATE TYPE tournament_system AS ENUM ('knockout', 'round_robin', 'swiss', 'double_elimination', 'groups_playoff');
CREATE TYPE participant_status AS ENUM ('registered', 'confirmed', 'checked_in', 'no_show', 'cancelled');

CREATE TYPE result_status AS ENUM ('pending', 'confirmed', 'disputed', 'admin_confirmed');

CREATE TYPE chat_type AS ENUM ('personal', 'community', 'event');

CREATE TYPE notification_type AS ENUM (
    'event_response', 'game_reminder_24h', 'game_reminder_1h',
    'result_confirm', 'community_news', 'new_message',
    'rating_change', 'new_badge', 'join_request',
    'join_approved', 'join_rejected', 'event_cancelled', 'spot_available'
);

CREATE TYPE court_surface AS ENUM ('hard', 'clay', 'carpet', 'grass', 'synthetic');
CREATE TYPE post_author_type AS ENUM ('user', 'community');

-- =====================================================
-- 1. USERS
-- =====================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone VARCHAR(20) UNIQUE NOT NULL,
    phone_verified BOOLEAN DEFAULT FALSE,

    -- Profile
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    gender gender_type,
    birth_year SMALLINT CHECK (birth_year >= 1940 AND birth_year <= 2015),
    city VARCHAR(100) DEFAULT 'Астана',
    district VARCHAR(100),
    avatar_url TEXT,
    bio TEXT,

    -- Tennis level
    ntrp_level DECIMAL(2,1) CHECK (ntrp_level >= 1.0 AND ntrp_level <= 7.0),
    level_label VARCHAR(50),
    quiz_completed BOOLEAN DEFAULT FALSE,

    -- Rating
    global_rating DECIMAL(7,2) DEFAULT 1000.00,
    global_games_count INT DEFAULT 0,

    -- Settings
    language VARCHAR(5) DEFAULT 'ru',
    pin_hash VARCHAR(255),
    push_token TEXT,
    platform_role platform_role DEFAULT 'player',

    -- Privacy
    profile_visibility VARCHAR(20) DEFAULT 'all',
    allow_messages_from VARCHAR(20) DEFAULT 'all',
    show_stats BOOLEAN DEFAULT TRUE,

    -- Notification settings
    notification_settings JSONB DEFAULT '{
        "event_response": true, "game_reminder_24h": true,
        "game_reminder_1h": true, "result_confirm": true,
        "community_news": true, "new_message": true,
        "rating_change": true, "new_badge": true,
        "quiet_hours_start": "23:00", "quiet_hours_end": "07:00"
    }'::jsonb,

    -- Status
    status user_status DEFAULT 'active',
    is_profile_complete BOOLEAN DEFAULT FALSE,
    last_active_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_status ON users(status) WHERE status = 'active';
CREATE INDEX idx_users_ntrp ON users(ntrp_level);
CREATE INDEX idx_users_city_district ON users(city, district);
CREATE INDEX idx_users_global_rating ON users(global_rating DESC);
CREATE INDEX idx_users_name_trgm ON users USING gin ((first_name || ' ' || last_name) gin_trgm_ops);
CREATE INDEX idx_users_last_active ON users(last_active_at DESC);

-- =====================================================
-- 2. AUTH / OTP
-- =====================================================

CREATE TABLE otp_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone VARCHAR(20) NOT NULL,
    code VARCHAR(6) NOT NULL,
    attempts INT DEFAULT 0,
    max_attempts INT DEFAULT 5,
    is_verified BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_otp_phone ON otp_sessions(phone);
CREATE INDEX idx_otp_expires ON otp_sessions(expires_at);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    device_info TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_refresh_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_hash ON refresh_tokens(token_hash);

-- =====================================================
-- 3. COURTS
-- =====================================================

CREATE TABLE courts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    district VARCHAR(100),
    latitude DECIMAL(10, 7),
    longitude DECIMAL(10, 7),
    total_courts SMALLINT DEFAULT 1,
    indoor_courts SMALLINT DEFAULT 0,
    outdoor_courts SMALLINT DEFAULT 0,
    surface court_surface,
    price_per_hour DECIMAL(10, 2),
    currency VARCHAR(3) DEFAULT 'KZT',
    phone VARCHAR(20),
    working_hours JSONB,
    photos JSONB DEFAULT '[]'::jsonb,
    community_id UUID,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_courts_location ON courts(latitude, longitude);
CREATE INDEX idx_courts_district ON courts(district);
CREATE INDEX idx_courts_active ON courts(is_active) WHERE is_active = TRUE;

-- =====================================================
-- 4. COMMUNITIES
-- =====================================================

CREATE TABLE communities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE,
    description TEXT,
    rules TEXT,
    community_type community_type NOT NULL,
    access_level community_access DEFAULT 'open',
    verification_status verification_status DEFAULT 'none',
    verified_at TIMESTAMPTZ,
    verification_documents JSONB,
    logo_url TEXT,
    banner_url TEXT,
    contact_phone VARCHAR(20),
    contact_email VARCHAR(255),
    social_links JSONB,
    address TEXT,
    district VARCHAR(100),
    rating_initial DECIMAL(7,2) DEFAULT 1000.00,
    rating_k_factor INT DEFAULT 32,
    rating_min_games INT DEFAULT 3,
    member_count INT DEFAULT 0,
    event_count INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_comm_type ON communities(community_type);
CREATE INDEX idx_comm_access ON communities(access_level);
CREATE INDEX idx_comm_verification ON communities(verification_status);
CREATE INDEX idx_comm_active ON communities(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_comm_slug ON communities(slug);
CREATE INDEX idx_comm_name_trgm ON communities USING gin (name gin_trgm_ops);

-- Add FK from courts to communities
ALTER TABLE courts ADD CONSTRAINT fk_courts_community
    FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE SET NULL;
CREATE INDEX idx_courts_community ON courts(community_id);

-- =====================================================
-- 5. COMMUNITY MEMBERS
-- =====================================================

CREATE TABLE community_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role community_role DEFAULT 'member',
    status member_status DEFAULT 'active',
    application_message TEXT,
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMPTZ,
    community_rating DECIMAL(7,2) DEFAULT 1000.00,
    community_games_count INT DEFAULT 0,
    community_wins INT DEFAULT 0,
    community_losses INT DEFAULT 0,
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(community_id, user_id)
);

CREATE INDEX idx_cm_community ON community_members(community_id);
CREATE INDEX idx_cm_user ON community_members(user_id);
CREATE INDEX idx_cm_status ON community_members(community_id, status);
CREATE INDEX idx_cm_role ON community_members(community_id, role);
CREATE INDEX idx_cm_rating ON community_members(community_id, community_rating DESC);

-- =====================================================
-- 6. EVENTS
-- =====================================================

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_type event_type NOT NULL,
    status event_status DEFAULT 'draft',
    community_id UUID REFERENCES communities(id) ON DELETE SET NULL,
    player_composition player_composition NOT NULL DEFAULT 'singles',
    match_format match_format DEFAULT 'best_of',
    match_format_details JSONB,
    tournament_system tournament_system,
    tournament_details JSONB,
    court_id UUID REFERENCES courts(id),
    location_name VARCHAR(255),
    location_address TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    max_participants INT,
    min_participants INT DEFAULT 2,
    current_participants INT DEFAULT 0,
    min_level DECIMAL(2,1),
    max_level DECIMAL(2,1),
    gender_restriction gender_type,
    min_age SMALLINT,
    max_age SMALLINT,
    registration_deadline TIMESTAMPTZ,
    is_paid BOOLEAN DEFAULT FALSE,
    price_amount DECIMAL(10,2),
    price_currency VARCHAR(3) DEFAULT 'KZT',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_events_type ON events(event_type);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_community ON events(community_id);
CREATE INDEX idx_events_start ON events(start_time);
CREATE INDEX idx_events_creator ON events(created_by);
CREATE INDEX idx_events_composition ON events(player_composition);
CREATE INDEX idx_events_active ON events(status, start_time)
    WHERE status IN ('published', 'registration_open', 'in_progress');

-- =====================================================
-- 7. EVENT PARTICIPANTS
-- =====================================================

CREATE TABLE event_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status participant_status DEFAULT 'registered',
    registered_at TIMESTAMPTZ DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ,
    partner_id UUID REFERENCES users(id),
    seed_number INT,
    UNIQUE(event_id, user_id)
);

CREATE INDEX idx_ep_event ON event_participants(event_id);
CREATE INDEX idx_ep_user ON event_participants(user_id);
CREATE INDEX idx_ep_status ON event_participants(event_id, status);

-- =====================================================
-- 8. MATCHES & RESULTS
-- =====================================================

CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID REFERENCES events(id) ON DELETE SET NULL,
    community_id UUID REFERENCES communities(id) ON DELETE SET NULL,

    player1_id UUID NOT NULL REFERENCES users(id),
    player2_id UUID NOT NULL REFERENCES users(id),
    player1_partner_id UUID REFERENCES users(id),
    player2_partner_id UUID REFERENCES users(id),

    composition player_composition NOT NULL DEFAULT 'singles',

    score JSONB,
    winner_id UUID REFERENCES users(id),

    result_status result_status DEFAULT 'pending',
    submitted_by UUID REFERENCES users(id),
    confirmed_by UUID REFERENCES users(id),
    submitted_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ,
    dispute_reason TEXT,

    player1_rating_before DECIMAL(7,2),
    player1_rating_after DECIMAL(7,2),
    player2_rating_before DECIMAL(7,2),
    player2_rating_after DECIMAL(7,2),

    round_name VARCHAR(50),
    round_number INT,
    court_number INT,
    scheduled_time TIMESTAMPTZ,

    played_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CHECK (player1_id != player2_id)
);

CREATE INDEX idx_matches_event ON matches(event_id);
CREATE INDEX idx_matches_community ON matches(community_id);
CREATE INDEX idx_matches_p1 ON matches(player1_id);
CREATE INDEX idx_matches_p2 ON matches(player2_id);
CREATE INDEX idx_matches_winner ON matches(winner_id);
CREATE INDEX idx_matches_status ON matches(result_status);
CREATE INDEX idx_matches_played ON matches(played_at DESC);

-- =====================================================
-- 9. CHATS & MESSAGES
-- =====================================================

CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_type chat_type NOT NULL,
    community_id UUID REFERENCES communities(id) ON DELETE CASCADE,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    user1_id UUID REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255),
    is_archived BOOLEAN DEFAULT FALSE,
    last_message_at TIMESTAMPTZ,
    last_message_preview TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CHECK (
        (chat_type = 'personal' AND user1_id IS NOT NULL AND user2_id IS NOT NULL AND user1_id < user2_id)
        OR (chat_type = 'community' AND community_id IS NOT NULL)
        OR (chat_type = 'event' AND event_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX idx_chats_personal ON chats(user1_id, user2_id) WHERE chat_type = 'personal';
CREATE UNIQUE INDEX idx_chats_community ON chats(community_id) WHERE chat_type = 'community';
CREATE UNIQUE INDEX idx_chats_event ON chats(event_id) WHERE chat_type = 'event';
CREATE INDEX idx_chats_type ON chats(chat_type);
CREATE INDEX idx_chats_last_msg ON chats(last_message_at DESC);

CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    reply_to_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_msg_chat ON messages(chat_id, created_at DESC);
CREATE INDEX idx_msg_sender ON messages(sender_id);

CREATE TABLE chat_read_status (
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_at TIMESTAMPTZ DEFAULT NOW(),
    is_muted BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (chat_id, user_id)
);

-- Add pinned message column to chats
ALTER TABLE chats ADD COLUMN pinned_message_id UUID REFERENCES messages(id) ON DELETE SET NULL;

-- =====================================================
-- 10. NOTIFICATIONS
-- =====================================================

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type notification_type NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    data JSONB,
    is_read BOOLEAN DEFAULT FALSE,
    is_pushed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_notif_user ON notifications(user_id, created_at DESC);
CREATE INDEX idx_notif_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_notif_expires ON notifications(expires_at) WHERE expires_at IS NOT NULL;

-- =====================================================
-- 11. FRIENDS
-- =====================================================

CREATE TABLE friends (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, friend_id),
    CHECK(user_id != friend_id)
);

CREATE INDEX idx_friends_user ON friends(user_id);
CREATE INDEX idx_friends_friend ON friends(friend_id);

-- =====================================================
-- 12. POSTS / FEED
-- =====================================================

CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_type post_author_type NOT NULL,
    author_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    author_community_id UUID REFERENCES communities(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    photos JSONB DEFAULT '[]'::jsonb,
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    is_match_result BOOLEAN DEFAULT FALSE,
    match_id UUID REFERENCES matches(id) ON DELETE SET NULL,
    is_published BOOLEAN DEFAULT TRUE,
    scheduled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CHECK (
        (author_type = 'user' AND author_user_id IS NOT NULL)
        OR (author_type = 'community' AND author_community_id IS NOT NULL)
    )
);

CREATE INDEX idx_posts_user ON posts(author_user_id, created_at DESC);
CREATE INDEX idx_posts_community ON posts(author_community_id, created_at DESC);
CREATE INDEX idx_posts_created ON posts(created_at DESC) WHERE is_published = TRUE;

CREATE TABLE post_likes (
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)
);

-- =====================================================
-- 13. BADGES / ACHIEVEMENTS
-- =====================================================

CREATE TABLE badge_definitions (
    id VARCHAR(50) PRIMARY KEY,
    name_ru VARCHAR(100) NOT NULL,
    name_kz VARCHAR(100),
    name_en VARCHAR(100),
    description_ru TEXT,
    description_kz TEXT,
    description_en TEXT,
    icon VARCHAR(10),
    condition_type VARCHAR(50) NOT NULL,
    condition_value INT NOT NULL,
    sort_order INT DEFAULT 0
);

CREATE TABLE user_badges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id VARCHAR(50) NOT NULL REFERENCES badge_definitions(id),
    earned_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, badge_id)
);

CREATE INDEX idx_ubadges_user ON user_badges(user_id);

-- =====================================================
-- 14. PLAYER STATS (denormalized cache)
-- =====================================================

CREATE TABLE player_stats_global (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_games INT DEFAULT 0,
    total_wins INT DEFAULT 0,
    total_losses INT DEFAULT 0,
    win_rate DECIMAL(5,2) DEFAULT 0.00,
    singles_games INT DEFAULT 0,
    singles_wins INT DEFAULT 0,
    doubles_games INT DEFAULT 0,
    doubles_wins INT DEFAULT 0,
    current_streak INT DEFAULT 0,
    best_streak INT DEFAULT 0,
    tournaments_played INT DEFAULT 0,
    last_game_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- =====================================================
-- 15. RATING HISTORY
-- =====================================================

CREATE TABLE rating_history (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    community_id UUID REFERENCES communities(id) ON DELETE CASCADE,
    rating_before DECIMAL(7,2) NOT NULL,
    rating_after DECIMAL(7,2) NOT NULL,
    change DECIMAL(7,2) NOT NULL,
    match_id UUID REFERENCES matches(id) ON DELETE SET NULL,
    reason VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_rh_user ON rating_history(user_id, created_at DESC);
CREATE INDEX idx_rh_user_comm ON rating_history(user_id, community_id, created_at DESC);

-- =====================================================
-- 16. AUDIT LOG
-- =====================================================

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50),
    entity_id UUID,
    details JSONB,
    ip_address INET,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_actor ON audit_logs(actor_id);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at DESC);

-- =====================================================
-- FUNCTIONS & TRIGGERS
-- =====================================================

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_communities_updated BEFORE UPDATE ON communities FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_cm_updated BEFORE UPDATE ON community_members FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_events_updated BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_matches_updated BEFORE UPDATE ON matches FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_chats_updated BEFORE UPDATE ON chats FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_posts_updated BEFORE UPDATE ON posts FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_courts_updated BEFORE UPDATE ON courts FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Auto-update community member_count
CREATE OR REPLACE FUNCTION update_community_member_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE communities SET member_count = (
        SELECT COUNT(*) FROM community_members
        WHERE community_id = COALESCE(NEW.community_id, OLD.community_id)
          AND status = 'active'
    ) WHERE id = COALESCE(NEW.community_id, OLD.community_id);
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_cm_count
AFTER INSERT OR UPDATE OR DELETE ON community_members
FOR EACH ROW EXECUTE FUNCTION update_community_member_count();

-- Auto-update event participant_count
CREATE OR REPLACE FUNCTION update_event_participant_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE events SET current_participants = (
        SELECT COUNT(*) FROM event_participants
        WHERE event_id = COALESCE(NEW.event_id, OLD.event_id)
          AND status IN ('registered', 'confirmed', 'checked_in')
    ) WHERE id = COALESCE(NEW.event_id, OLD.event_id);
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_ep_count
AFTER INSERT OR UPDATE OR DELETE ON event_participants
FOR EACH ROW EXECUTE FUNCTION update_event_participant_count();

-- =====================================================
-- VIEWS
-- =====================================================

CREATE OR REPLACE VIEW v_community_leaderboard AS
SELECT
    cm.community_id,
    cm.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.ntrp_level,
    cm.community_rating,
    cm.community_games_count,
    cm.community_wins,
    cm.community_losses,
    CASE WHEN cm.community_games_count > 0
        THEN ROUND(cm.community_wins::decimal / cm.community_games_count * 100, 1)
        ELSE 0
    END as win_rate,
    ROW_NUMBER() OVER (
        PARTITION BY cm.community_id
        ORDER BY cm.community_rating DESC
    ) as rank
FROM community_members cm
JOIN users u ON cm.user_id = u.id
WHERE cm.status = 'active' AND u.status = 'active';

CREATE OR REPLACE VIEW v_active_events AS
SELECT
    e.*,
    c.name as community_name,
    c.logo_url as community_logo,
    c.community_type,
    u.first_name as creator_first_name,
    u.last_name as creator_last_name,
    ct.name as court_name,
    (e.max_participants - e.current_participants) as available_spots
FROM events e
LEFT JOIN communities c ON e.community_id = c.id
JOIN users u ON e.created_by = u.id
LEFT JOIN courts ct ON e.court_id = ct.id
WHERE e.status IN ('published', 'registration_open', 'in_progress')
ORDER BY e.start_time;

-- =====================================================
-- SEED DATA
-- =====================================================

-- Badge definitions
INSERT INTO badge_definitions (id, name_ru, name_en, icon, condition_type, condition_value, sort_order) VALUES
('first_win',     'Первая победа',    'First Win',      '🏅', 'wins',          1,   1),
('ten_wins',      'Десятка',          'Ten Wins',       '🎖', 'wins',          10,  2),
('fifty_games',   'Полтинник',        'Fifty Games',    '⭐', 'games',         50,  3),
('win_streak_5',  'Серия',            'Win Streak',     '🔥', 'streak',        5,   4),
('tournaments_5', 'Турнирный боец',   'Tournament Pro', '🏆', 'tournaments',   5,   5),
('level_up',      'Уровень вверх',    'Level Up',       '👑', 'level_up',      1,   6),
('social_3',      'Социальный',       'Social',         '🤝', 'communities',   3,   7),
('veteran_100',   'Ветеран',          'Veteran',        '🎾', 'games',         100, 8),
('friends_10',    'Авторитет',        'Popular',        '👥', 'friends',       10,  9),
('weekly_4',      'Регулярный',       'Regular',        '📅', 'weekly_streak', 4,   10);

-- Courts of Astana
INSERT INTO courts (name, address, district, latitude, longitude, total_courts, indoor_courts, outdoor_courts, surface, price_per_hour, phone) VALUES
('NTC Astana',          'Кабанбай батыра, 42',    'Есильский',     51.1282, 71.4307, 6, 4, 2, 'hard',   3000, '+77172123456'),
('Tennis Club Astana',  'Мәңгілік Ел, 15',        'Есильский',     51.1320, 71.4185, 4, 4, 0, 'hard',   4000, '+77172234567'),
('Mega Tennis',         'Қабанбай батыра, 62/5',   'Сарыаркинский', 51.1450, 71.4102, 3, 3, 0, 'carpet', 3500, '+77172345678'),
('Президентский клуб',  'Бейбітшілік, 10',        'Алматинский',   51.1180, 71.4290, 8, 4, 4, 'hard',   5000, '+77172456789');
