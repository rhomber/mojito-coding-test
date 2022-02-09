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
    created_at DATE NOT NULL,
    updated_at DATE,
    deleted_at DATE
);

CREATE TABLE "auction_lot_bid"
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    auction_lot_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    type TEXT NOT NULL,
    bid INTEGER NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE,
    deleted_at DATE,
    FOREIGN KEY(auction_lot_id) REFERENCES auction_lot(id),
    FOREIGN KEY(user_id) REFERENCES user(id)
);

-- I would usually put a unique constraint around active but I am unfamiliar with sqlite3 to know if it would
-- work correctly (ignoring nulls) so I have avoided it for now.
CREATE TABLE "auction_lot_bid_max"
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    auction_lot_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    max_bid INTEGER NOT NULL,
    active BOOL,
    created_at DATE NOT NULL,
    updated_at DATE,
    deleted_at DATE,
    FOREIGN KEY(auction_lot_id) REFERENCES auction_lot(id),
    FOREIGN KEY(user_id) REFERENCES user(id)
);

-- +goose Down
DROP TABLE "user";
DROP TABLE "auction_lot";
DROP TABLE "auction_lot_bid";
DROP TABLE "auction_lot_bid_max";