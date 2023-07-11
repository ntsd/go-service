CREATE TABLE IF NOT EXISTS "clients" (
  "id" text NOT NULL UNIQUE PRIMARY KEY,
  "secret" text NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
