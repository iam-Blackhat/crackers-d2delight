-- 1. Add the role_id column (nullable first to avoid issues)
ALTER TABLE users
ADD COLUMN role_id UUID;

-- 2. Add foreign key constraint to roles table
ALTER TABLE users
ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id);

-- 3. Enforce NOT NULL once you confirm the app always sets role_id
ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;
