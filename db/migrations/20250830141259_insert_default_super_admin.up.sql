-- super_admin_seed.sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO users (name, email, phone, password, role_id, created_at, updated_at)
VALUES (
  'Super Admin',
  'super@example.com',
  '0000000000',
  crypt('SuperAdmin123!', gen_salt('bf')),
  (SELECT id FROM roles WHERE name = 'SUPER ADMIN'),
  NOW(),
  NOW()
)
ON CONFLICT (email) DO NOTHING;

