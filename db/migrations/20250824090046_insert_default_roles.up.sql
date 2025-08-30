-- Up migration: insert default roles
INSERT INTO roles (id, name, created_at, updated_at)
VALUES
  (uuid_generate_v4(), 'SUPER ADMIN', NOW(), NOW()),
  (uuid_generate_v4(), 'ADMIN', NOW(), NOW()),
  (uuid_generate_v4(), 'CUSTOMER', NOW(), NOW());
