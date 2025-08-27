-- 1. Rename the table
ALTER TABLE customers RENAME TO customer_profiles;

-- 2. Drop the old username column
ALTER TABLE customer_profiles DROP COLUMN username;

-- 3. Add new user_id column as BIGINT
ALTER TABLE customer_profiles
ADD COLUMN user_id BIGINT NOT NULL;

-- 4. Add foreign key constraint
ALTER TABLE customer_profiles
ADD CONSTRAINT fk_customer_user
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
