-- 1. Drop foreign key constraint if it exists
ALTER TABLE customers DROP CONSTRAINT IF EXISTS customers_user_id_fkey;

-- 2. Drop the user_id column
ALTER TABLE customers DROP COLUMN IF EXISTS user_id;

-- 3. Add username column
ALTER TABLE customers ADD COLUMN IF NOT EXISTS username VARCHAR(255) NOT NULL DEFAULT '';

-- 4. Ensure deleted_at column is nullable (soft deletes)
ALTER TABLE customers ALTER COLUMN deleted_at DROP NOT NULL;
