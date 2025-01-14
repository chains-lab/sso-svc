CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE role_type AS ENUM (
    'admin',
    'user',
    'verify_user'
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    role role_type DEFAULT 'user' NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    client TEXT NOT NULL,
    IP_first TEXT NOT NULL,
    IP_last TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    last_used TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_account_email ON accounts(email);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_last_used ON sessions(last_used);
