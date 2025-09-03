CREATE TABLE phone_otps (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL,
    purpose VARCHAR(50) NOT NULL,
    code_hash TEXT NOT NULL,
    meta JSONB,
    attempts INT DEFAULT 0,
    max_attempts INT DEFAULT 5,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_phone_otps_phone_purpose ON phone_otps(phone, purpose);
