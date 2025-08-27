-- Make phone unique
ALTER TABLE customers
ADD CONSTRAINT customers_phone_unique UNIQUE (phone);

-- Add email column with unique constraint
ALTER TABLE customers
ADD COLUMN email VARCHAR(255)
