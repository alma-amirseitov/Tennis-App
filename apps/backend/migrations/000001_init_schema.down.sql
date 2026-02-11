-- =====================================================
-- Reverse migration: 000001_init_schema
-- Drop in reverse dependency order
-- =====================================================

-- Drop views
DROP VIEW IF EXISTS v_active_events;
DROP VIEW IF EXISTS v_community_leaderboard;

-- Drop triggers
DROP TRIGGER IF EXISTS trg_ep_count ON event_participants;
DROP TRIGGER IF EXISTS trg_cm_count ON community_members;
DROP TRIGGER IF EXISTS trg_courts_updated ON courts;
DROP TRIGGER IF EXISTS trg_posts_updated ON posts;
DROP TRIGGER IF EXISTS trg_chats_updated ON chats;
DROP TRIGGER IF EXISTS trg_matches_updated ON matches;
DROP TRIGGER IF EXISTS trg_events_updated ON events;
DROP TRIGGER IF EXISTS trg_cm_updated ON community_members;
DROP TRIGGER IF EXISTS trg_communities_updated ON communities;
DROP TRIGGER IF EXISTS trg_users_updated ON users;

-- Drop functions
DROP FUNCTION IF EXISTS update_event_participant_count();
DROP FUNCTION IF EXISTS update_community_member_count();
DROP FUNCTION IF EXISTS update_updated_at();

-- Drop tables (reverse dependency order)
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS rating_history CASCADE;
DROP TABLE IF EXISTS player_stats_global CASCADE;
DROP TABLE IF EXISTS user_badges CASCADE;
DROP TABLE IF EXISTS badge_definitions CASCADE;
DROP TABLE IF EXISTS post_likes CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS friends CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS chat_read_status CASCADE;
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS chats CASCADE;
DROP TABLE IF EXISTS matches CASCADE;
DROP TABLE IF EXISTS event_participants CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS community_members CASCADE;
DROP TABLE IF EXISTS courts CASCADE;
DROP TABLE IF EXISTS communities CASCADE;
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS otp_sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop enums
DROP TYPE IF EXISTS post_author_type;
DROP TYPE IF EXISTS court_surface;
DROP TYPE IF EXISTS notification_type;
DROP TYPE IF EXISTS chat_type;
DROP TYPE IF EXISTS result_status;
DROP TYPE IF EXISTS participant_status;
DROP TYPE IF EXISTS tournament_system;
DROP TYPE IF EXISTS match_format;
DROP TYPE IF EXISTS player_composition;
DROP TYPE IF EXISTS event_status;
DROP TYPE IF EXISTS event_type;
DROP TYPE IF EXISTS member_status;
DROP TYPE IF EXISTS community_role;
DROP TYPE IF EXISTS verification_status;
DROP TYPE IF EXISTS community_access;
DROP TYPE IF EXISTS community_type;
DROP TYPE IF EXISTS platform_role;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS gender_type;

-- Drop extensions
DROP EXTENSION IF EXISTS "pg_trgm";
DROP EXTENSION IF EXISTS "uuid-ossp";
