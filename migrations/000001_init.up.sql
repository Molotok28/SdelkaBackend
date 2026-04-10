DROP SCHEMA IF EXISTS sdelka CASCADE;
CREATE SCHEMA sdelka;

CREATE TABLE sdelka.users
(
    id           SERIAL PRIMARY KEY,
    version      BIGINT       NOT NULL DEFAULT 1,
    name         VARCHAR(100) NOT NULL CHECK (char_length(name) between 3 and 100),
    surname         VARCHAR(100) NOT NULL CHECK (char_length(name) between 3 and 100),
    phone_number VARCHAR(15) CHECK (
        phone_number ~ '^\+?[1-9]\d{1,14}$' AND
        char_length (phone_number) between 10 and 15
)
    );

CREATE TABLE sdelka.ads
(
    id           SERIAL PRIMARY KEY,
    version      BIGINT       NOT NULL DEFAULT 1,
    title        VARCHAR(200) NOT NULL CHECK (char_length(title) between 5 and 200),
    description  TEXT         NOT NULL CHECK (char_length(description) between 10 and 1000),
    price        NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    user_id      INTEGER      NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES sdelka.users(id) ON DELETE CASCADE
)