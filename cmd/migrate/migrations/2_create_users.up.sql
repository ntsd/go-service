CREATE TABLE IF NOT EXISTS "users" (
  "id" text NOT NULL UNIQUE PRIMARY KEY,
  "email" text NOT NULL UNIQUE,
  "name" text NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- create GIN pg_trgm index
CREATE INDEX "users_name_trgm_idx" ON "users" USING GIN ("name" gin_trgm_ops);
