CREATE TABLE IF NOT EXISTS mutable.traffic (
    city_id INT NOT NULL REFERENCES mutable.cities(id) ON DELETE CASCADE,
    employment_id INT NOT NULL REFERENCES mutable.employments(id) ON DELETE CASCADE,
    q SMALLINT NOT NULL CHECK (q >= 0 AND q <= 127),
    r SMALLINT NOT NULL CHECK (r >= 0 AND r <= 127),
    workers_anchor_q SMALLINT NOT NULL CHECK (workers_anchor_q >= 0 AND workers_anchor_q <= 127),
    workers_anchor_r SMALLINT NOT NULL CHECK (workers_anchor_r >= 0 AND workers_anchor_r <= 127),
    foot_value INT NOT NULL DEFAULT 0,
    car_value INT NOT NULL DEFAULT 0,
    bus_value INT NOT NULL DEFAULT 0,
    bicycle_value INT NOT NULL DEFAULT 0,
    train_value INT NOT NULL DEFAULT 0,
    subway_value INT NOT NULL DEFAULT 0,

    PRIMARY KEY (city_id, employment_id, q, r)
);

CREATE INDEX IF NOT EXISTS index_traffic_city_position
    ON mutable.traffic (city_id, q, r);

CREATE INDEX IF NOT EXISTS index_traffic_workers
    ON mutable.traffic (city_id, workers_anchor_q, workers_anchor_r);

---- create above / drop below ----
DROP TABLE IF EXISTS mutable.traffic;
