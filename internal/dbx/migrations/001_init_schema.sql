-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id               UUID         PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    role             VARCHAR(255) DEFAULT 'user' NOT NULL,
    email            TEXT         NOT NULL UNIQUE,
    email_verified   BOOLEAN      DEFAULT FALSE NOT NULL,
    email_updated_at TIMESTAMP    NOT NULL DEFAULT now(),
    created_at       TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE TABLE user_passwords (
    user_id        UUID         PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password_hash  TEXT         NOT NULL,
    updated_at     TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id         UUID      PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id    UUID      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT      NOT NULL,
    client     TEXT      NOT NULL,
    ip         TEXT      NOT NULL,
    last_used  TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_email ON users(email);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_last_used ON sessions(last_used);

-- +migrate Down
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP EXTENSION IF EXISTS "uuid-ossp";
