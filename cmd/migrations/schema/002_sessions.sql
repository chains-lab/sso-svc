-- +migrate Up

CREATE TABLE sessions (
    id         UUID      PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id    UUID      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT      NOT NULL,
    last_used  TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS sessions CASCADE;