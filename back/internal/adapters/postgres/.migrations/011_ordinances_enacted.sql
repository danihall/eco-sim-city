CREATE TABLE IF NOT EXISTS mutable.ordinances_enacted (
    city_id INT NOT NULL REFERENCES mutable.cities(id) ON DELETE CASCADE,
    ordinance_id SMALLINT NOT NULL REFERENCES static.ordinances(id),
    PRIMARY KEY (city_id, ordinance_id)
);

---- create above / drop below ----
DROP TABLE IF EXISTS mutable.ordinances_enacted;
