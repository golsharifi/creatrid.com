DROP TRIGGER IF EXISTS accounts_updated_at ON accounts;
DROP TRIGGER IF EXISTS users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at();
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
