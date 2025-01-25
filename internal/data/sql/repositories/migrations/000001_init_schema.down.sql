-- Удаление таблиц
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;

-- Удаление ENUM-типов
DROP TYPE IF EXISTS failure_reason;
DROP TYPE IF EXISTS role_type;

DROP EXTENSION IF EXISTS "uuid-ossp";
