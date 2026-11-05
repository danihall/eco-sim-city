CREATE TABLE IF NOT EXISTS static.ordinances (
    id SMALLINT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

---- create above / drop below ----
DROP TABLE IF EXISTS static.ordinances;
