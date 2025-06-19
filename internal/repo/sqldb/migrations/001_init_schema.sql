-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    role VARCHAR(255) DEFAULT 'user' NOT NULL,
    verified BOOLEAN DEFAULT FALSE NOT NULL,
    subscription UUID,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    client TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    last_used TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_email ON users(email);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_last_used ON sessions(last_used);

-- +migrate Down
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP EXTENSION IF EXISTS "uuid-ossp";
