-- +goose Up
CREATE TABLE "user"
(
  id integer PRIMARY KEY AUTOINCREMENT
);

-- +goose Down
DROP TABLE "user";