-- Migration 0001: initial schema
CREATE TABLE users (
    user_id       SERIAL PRIMARY KEY,
    email         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    username      TEXT        NOT NULL,
    register_date TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE urls (
    url_id        SERIAL PRIMARY KEY,
    original_url  TEXT        NOT NULL,
    short_code    TEXT        NOT NULL UNIQUE,
    custom_alias  TEXT,
    creation_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    user_id       INT         NOT NULL REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE visits (
    click_id   SERIAL PRIMARY KEY,
    url_id     INT         NOT NULL REFERENCES urls(url_id) ON DELETE CASCADE,
    click_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    referrer   TEXT,
    ip_address TEXT        NOT NULL,
    country    TEXT
);

-- Migration 0002: remove username column
ALTER TABLE users
    DROP COLUMN username;

-- Migration 0003: make custom_alias NOT NULL
ALTER TABLE urls
    ALTER COLUMN custom_alias SET NOT NULL;
