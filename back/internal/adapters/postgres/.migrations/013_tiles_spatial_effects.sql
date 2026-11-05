CREATE TABLE IF NOT EXISTS mutable.tiles_spatial_effects (
    city_id INT NOT NULL REFERENCES mutable.cities(id) ON DELETE CASCADE,
    q SMALLINT NOT NULL CHECK (q >= 0 AND q <= 127),
    r SMALLINT NOT NULL CHECK (r >= 0 AND r <= 127),
    ghg_1 SMALLINT NOT NULL,
    land_value SMALLINT NOT NULL,
    crime_generic SMALLINT NOT NULL DEFAULT 0,
    crime_vss SMALLINT NOT NULL DEFAULT 0,
    crime_drug SMALLINT NOT NULL DEFAULT 0,
    PRIMARY KEY (city_id, q, r)
);

---- create above / drop below ----
DROP TABLE IF EXISTS mutable.tiles_spatial_effects;
