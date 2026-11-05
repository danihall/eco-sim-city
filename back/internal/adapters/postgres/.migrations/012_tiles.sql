CREATE TABLE IF NOT EXISTS mutable.tiles (
    city_id INT NOT NULL REFERENCES mutable.cities(id) ON DELETE CASCADE,
    super_kind_id SMALLINT NOT NULL REFERENCES static.super_kinds(id),
    kind_id SMALLINT NOT NULL REFERENCES static.kinds (id),
    artifact_id SMALLINT NOT NULL REFERENCES static.artifacts(id),
    artifact_label TEXT NOT NULL,
    artifact_status SMALLINT NOT NULL CHECK (artifact_status >= 1 AND artifact_status <= 3),
    q SMALLINT NOT NULL CHECK (q >= 0 AND q <= 127),
    r SMALLINT NOT NULL CHECK (r >= 0 AND r <= 127),
    anchor_q SMALLINT NOT NULL CHECK (anchor_q >= 0 AND anchor_q <= 127),
    anchor_r SMALLINT NOT NULL CHECK (anchor_r >= 0 AND anchor_r <= 127),
    has_all_citizen_allocated BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY (city_id, q, r)
);

CREATE INDEX index_tiles_city_artifact ON mutable.tiles (city_id, artifact_id);

COMMENT ON COLUMN mutable.tiles.artifact_status IS '1: Step_1_AwaitingAssignment, 2: Step_2_AssignedAwaitingDelivery, 3: Step_3_DeliveredOperational';

---- create above / drop below ----
DROP TABLE IF EXISTS mutable.tiles;
