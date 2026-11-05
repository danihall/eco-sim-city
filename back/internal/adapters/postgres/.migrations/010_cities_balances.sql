CREATE TABLE IF NOT EXISTS mutable.cities_balances (
    city_id INT NOT NULL REFERENCES mutable.cities(id) ON DELETE CASCADE,
    revenues INT NOT NULL DEFAULT 0,
    expenses INT NOT NULL DEFAULT 0,
    PRIMARY KEY (city_id)
);

---- create above / drop below ----
DROP TABLE IF EXISTS mutable.cities_balances;
