-- Down migration: remove the default roles
DELETE FROM roles
WHERE name IN ('admin', 'staff', 'customer');
