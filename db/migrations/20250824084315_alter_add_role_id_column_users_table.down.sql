-- Down Migration

-- 1. Drop the foreign key constraint
ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_role;

-- 2. Drop the role_id column
ALTER TABLE users
DROP COLUMN IF EXISTS role_id;
