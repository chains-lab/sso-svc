CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE role_type AS ENUM (
    'admin',
    'user',
    'verify_user'
);

CREATE TYPE failure_reason AS ENUM (
    'invalid_password',
    'account_locked',
    'expired_token',
    'invalid_device_id',
    'invalid_refresh_token',
    'invalid_device_factory_id',
    'invalid_user_id',
    'too_many_attempts',
    'no_access',
    'success'
);

CREATE TYPE operation_type AS ENUM (
    'login',
    'refresh_token',
    'change_username',
    'change_password',
    'reset_password',
    'change_email',
    'delete_account',
    'delete_session',
    'terminate_session'
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
    device_name TEXT NOT NULL,
    client TEXT NOT NULL,
    IP TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    last_used TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_account_email ON accounts(email);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_last_used ON sessions(last_used);
CREATE INDEX idx_operations_user_id ON operations(user_id);
CREATE INDEX idx_operations_operation ON operations(operation);
CREATE INDEX idx_operations_created_at ON operations(created_at);
