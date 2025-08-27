-- 1. Drop the username column
ALTER TABLE customers DROP COLUMN IF EXISTS username;

-- 2. Add back the user_id column
ALTER TABLE customers ADD COLUMN IF NOT EXISTS user_id BIGINT;

-- 3. Re-add the foreign key constraint (assuming users.id is BIGINT PK)
ALTER TABLE customers
  ADD CONSTRAINT customers_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);

-- 4. (Optional) Ensure deleted_at column allows NOT NULL again if original was strict
-- ALTER TABLE customers ALTER COLUMN deleted_at SET NOT NULL;
