-- Migration Down

-- Drop triggers
DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_todos_due_date;
DROP INDEX IF EXISTS idx_todos_priority;
DROP INDEX IF EXISTS idx_todos_status;
DROP INDEX IF EXISTS idx_todos_user_id;

-- Drop tables
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp";
