-- +goose Up
CREATE TABLE "user"
(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  -- in real life (would be a secure hash):
  --password_hash TEXT NOT NULL
  created_at DATE NOT NULL,
  updated_at DATE,
  deleted_at DATE
);

CREATE TABLE "auction_lot"
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    start_time DATE NOT NULL,
    end_time DATE NOT NULL,
    -- in real life (would be a secure hash):
    --password_hash TEXT NOT NULL
    created_at DATE NOT NULL,
    updated_at DATE,
    deleted_at DATE
);


-- +goose Down
DROP TABLE "user";
DROP TABLE "auction_lot";