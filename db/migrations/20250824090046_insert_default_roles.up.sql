-- Up migration: insert default roles
INSERT INTO roles (id, name, created_at, updated_at)
VALUES
  (uuid_generate_v4(), 'super admin', NOW(), NOW()),
  (uuid_generate_v4(), 'admin', NOW(), NOW()),
  (uuid_generate_v4(), 'customer', NOW(), NOW());
