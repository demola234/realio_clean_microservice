CREATE TABLE "users" (
    "email" varchar UNIQUE NOT NULL,
    "hashed_password" varchar NOT NULL DEFAULT '',
    "full_name" varchar NOT NULL DEFAULT '',
    "role" varchar NOT NULL DEFAULT '',
    "created_at" timestamptz NOT NULL DEFAULT(now())
);