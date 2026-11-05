-- "static" schema encompasses authoritative data, thus they're never updated by anyone (except me)
CREATE SCHEMA IF NOT EXISTS static;

---- create above / drop below ----
DROP SCHEMA IF EXISTS static CASCADE;
