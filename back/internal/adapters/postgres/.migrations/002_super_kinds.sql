CREATE TABLE IF NOT EXISTS static.super_kinds (
    id SMALLINT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

---- create above / drop below ----
DROP TABLE IF EXISTS static.super_kinds;
