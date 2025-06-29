CREATE TABLE "sessions" (
    "id" uuid NOT NULL PRIMARY KEY,
    "username" VARCHAR NOT NULL,
    "refresh_token" VARCHAR NOT NULL,
    "user_agent" VARCHAR NOT NULL,
    "ip_address" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    "expires_at" TIMESTAMPTZ NOT NULL,
    "is_blocked" BOOLEAN NOT NULL DEFAULT false
);

-- Create foreign key for username to users table
ALTER TABLE "sessions" ADD CONSTRAINT "sessions_username_fk" FOREIGN KEY ("username") REFERENCES "users" ("username");