-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "user_role" AS ENUM (
    'super_user',
    'admin',
    'moderator',
    'user'
);

CREATE TYPE "user_status" AS ENUM (
    'active',
    'blocked'
);

CREATE TABLE users (
    id         UUID        PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    role       user_role   DEFAULT 'user'   NOT NULL,
    status     user_status DEFAULT 'active' NOT NULL ,
    created_at TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE TABLE user_passwords (
    user_id       UUID      PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT      NOT NULL,
    updated_at    TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE user_emails (
    user_id  UUID    PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email    TEXT    UNIQUE NOT NULL,
    verified BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE sessions (
    id         UUID      PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id    UUID      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT      NOT NULL,
    last_used  TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS user_passwords CASCADE;
DROP TABLE IF EXISTS user_emails CASCADE;

DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS user_status;

DROP EXTENSION IF EXISTS "uuid-ossp";
