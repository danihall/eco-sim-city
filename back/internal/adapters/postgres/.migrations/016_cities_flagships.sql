CREATE TABLE IF NOT EXISTS mutable.cities_flagships (
    city_id INT NOT NULL REFERENCES mutable.cities(id) ON DELETE CASCADE,
    artifact_id SMALLINT NOT NULL REFERENCES static.artifacts(id),
    PRIMARY KEY (city_id, artifact_id)
);

---- create above / drop below ----
DROP TABLE IF EXISTS mutable.cities_flagships;
