-- Active: 1727625324173@@127.0.0.1@5432@gopaydb
CREATE TYPE "Role" AS ENUM('COMMON', 'MERCHANT');

CREATE TABLE IF NOT EXISTS users (
    "id" UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "document" VARCHAR(255) UNIQUE NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password_hash" VARCHAR(255) NOT NULL,
    "balance" DECIMAL NOT NULL DEFAULT 0,
    "role" "Role" NOT NULL DEFAULT 'COMMON',
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);
---- create above / drop below ----
DROP TABLE IF EXISTS users;