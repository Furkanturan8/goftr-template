-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON tokens;
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS token_blacklist;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS users;

-- Drop types
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role; 