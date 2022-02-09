-- +goose Up
CREATE TABLE "user"
(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  -- in real life (would be a secure hash):
  --password_hash TEXT NOT NULL
  created_at TEXT NOT NULL,
  updated_at TEXT,
  deleted_at TEXT
);

-- +goose Down
DROP TABLE "user";