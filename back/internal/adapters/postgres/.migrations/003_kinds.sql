CREATE TABLE IF NOT EXISTS static.kinds (
    id SMALLINT PRIMARY KEY,
    super_kind_id SMALLINT NOT NULL REFERENCES static.super_kinds(id),
    name TEXT NOT NULL UNIQUE
);

---- create above / drop below ----
DROP TABLE IF EXISTS static.kinds;
