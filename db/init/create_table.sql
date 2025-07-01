CREATE TABLE sessions (
    session_id SERIAL PRIMARY KEY,
    user_guid UUID NOT NULL,
    user_agent TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);