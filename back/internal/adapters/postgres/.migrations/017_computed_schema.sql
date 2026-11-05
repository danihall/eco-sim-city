-- all tables within "computed" schema are UNLOGGED tables,
-- they must be re-populated from "mutable" schema's data on each database boot (in theory never)
CREATE SCHEMA IF NOT EXISTS computed;

---- create above / drop below ----
DROP SCHEMA IF EXISTS computed CASCADE;
