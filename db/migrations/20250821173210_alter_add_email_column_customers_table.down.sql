-- Drop unique constraint on phone
ALTER TABLE customers
DROP CONSTRAINT IF EXISTS customers_phone_unique;

-- Drop email column (this also removes its unique constraint)
ALTER TABLE customers
DROP COLUMN IF EXISTS email;
