CREATE TABLE IF NOT EXISTS mutable.cities_budget (
    city_id INT NOT NULL REFERENCES mutable.cities(id) ON DELETE CASCADE,
    kind_id SMALLINT NOT NULL REFERENCES static.kinds(id),
    tax_rate SMALLINT NOT NULL,
    funding_rate SMALLINT NOT NULL,
    PRIMARY KEY (city_id, kind_id)
);

---- create above / drop below ----
DROP TABLE IF EXISTS mutable.cities_budget;
