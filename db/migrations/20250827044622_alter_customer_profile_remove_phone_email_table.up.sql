ALTER TABLE customer_profiles
DROP COLUMN email,
DROP COLUMN phone,
DROP COLUMN address,
ADD COLUMN addresses JSONB DEFAULT '[]'::jsonb;

