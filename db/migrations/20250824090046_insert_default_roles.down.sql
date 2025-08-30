-- Down migration: remove the default roles
DELETE FROM roles
WHERE name IN ('SUPER ADMIN', 'ADMIN', 'CUSTOMER');
