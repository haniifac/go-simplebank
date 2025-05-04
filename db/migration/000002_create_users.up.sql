CREATE TABLE users(
    username VARCHAR PRIMARY KEY,
    hashed_password VARCHAR NOT NULL,
    fullname VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    password_updated_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00+00:00'
);

ALTER TABLE accounts ADD CONSTRAINT "fk_account_owner" FOREIGN KEY (owner) REFERENCES users (username);
ALTER TABLE accounts ADD CONSTRAINT "unique_owner_currency" UNIQUE (owner, currency);

-- CREATE INDEX "users_username_idx" ON "users" ("username");