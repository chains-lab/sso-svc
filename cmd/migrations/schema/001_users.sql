-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "user_role" AS ENUM (
    'admin',
    'moderator',
    'user'
);

CREATE TYPE "user_status" AS ENUM (
    'active',
    'blocked'
);

CREATE TABLE users (
    id                  UUID        PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    role                user_role   DEFAULT 'user'   NOT NULL,
    status              user_status DEFAULT 'active' NOT NULL,

    password_hash       TEXT        NOT NULL,
    password_updated_at TIMESTAMP   NOT NULL DEFAULT now(),

    email              TEXT        UNIQUE NOT NULL,
    email_verified     BOOLEAN     DEFAULT FALSE NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS users CASCADE;

DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS user_status;

DROP EXTENSION IF EXISTS "uuid-ossp";
