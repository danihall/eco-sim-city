-- name: BuildArtifact :many
UPDATE
    mutable.tiles AS _tile
SET
    super_kind_id = @super_kind_id,
    kind_id = @kind_id,
    artifact_id = @artifact_id,
    artifact_label = @artifact_label,
    artifact_status = @artifact_status,
    anchor_q = _params.anchor_q,
    anchor_r = _params.anchor_r
FROM ROWS FROM(
    UNNEST(@qs::SMALLINT[]),
    UNNEST(@rs::SMALLINT[]),
    UNNEST(@anchor_qs::SMALLINT[]),
    UNNEST(@anchor_rs::SMALLINT[])
) AS _params(q, r, anchor_q, anchor_r)
WHERE
    _tile.city_id = @city_id
    AND _tile.q = _params.q
    AND _tile.r = _params.r
RETURNING
    _tile.super_kind_id,
    _tile.q,
    _tile.r,
    _tile.anchor_q,
    _tile.anchor_r,
    _tile.artifact_id,
    _tile.artifact_status;

-- name: BuildMultipleArtifacts :many
UPDATE
    mutable.tiles AS _tile
SET
    super_kind_id = _params.super_kind_id,
    kind_id = _params.kind_id,
    artifact_id = _params.artifact_id,
    artifact_label = _params.artifact_label,
    artifact_status = _params.artifact_status,
    anchor_q = _params.anchor_q,
    anchor_r = _params.anchor_r
FROM ROWS FROM(
    UNNEST(@qs::SMALLINT[]),
    UNNEST(@rs::SMALLINT[]),
    UNNEST(@super_kind_ids::SMALLINT[]),
    UNNEST(@kind_ids::SMALLINT[]),
    UNNEST(@artifact_ids::SMALLINT[]),
    UNNEST(@artifact_labels::TEXT[]),
    UNNEST(@artifact_statuses::SMALLINT[]),
    UNNEST(@anchor_qs::SMALLINT[]),
    UNNEST(@anchor_rs::SMALLINT[])
) AS _params(
    q,
    r,
    super_kind_id,
    kind_id,
    artifact_id,
    artifact_label,
    artifact_status,
    anchor_q,
    anchor_r
)
WHERE
    _tile.city_id = @city_id
    AND _tile.q = _params.q
    AND _tile.r = _params.r
RETURNING _tile.artifact_status;

-- name: BuildPathway :many
UPDATE
    mutable.tiles AS _tile
SET
    super_kind_id = @super_kind_id,
    kind_id = @kind_id,
    artifact_id = @artifact_id,
    artifact_label = @artifact_label,
    artifact_status = @completed_artifact_status,
    anchor_q = _params.anchor_q,
    anchor_r = _params.anchor_r
FROM ROWS FROM(
    UNNEST(@qs::SMALLINT[]),
    UNNEST(@rs::SMALLINT[]),
    UNNEST(@anchor_qs::SMALLINT[]),
    UNNEST(@anchor_rs::SMALLINT[])
) AS _params(q, r, anchor_q, anchor_r)
WHERE
    _tile.city_id = @city_id
    AND _tile.q = _params.q
    AND _tile.r = _params.r
RETURNING
    _tile.super_kind_id,
    _tile.q,
    _tile.r,
    _tile.anchor_q,
    _tile.anchor_r,
    _tile.artifact_id,
    _tile.artifact_label,
    _tile.artifact_status;

-- name: CompleteArtifactsFromChunk :many
UPDATE mutable.tiles
SET
    artifact_status = @new_artifact_status
WHERE
    city_id = @city_id
    AND artifact_status = @current_artifact_status
    AND q BETWEEN @q_top_left AND @q_bottom_right
    AND r BETWEEN @r_top_left AND @r_bottom_right
RETURNING
    q,
    r,
    anchor_q,
    anchor_r,
    artifact_id,
    artifact_status;

-- name: VacateMultipleArtifacts :many
UPDATE
    mutable.tiles AS _tile
SET
    artifact_label = @artifact_label,
    artifact_status = @artifact_status,
    has_all_citizen_allocated = @has_all_citizen_allocated
FROM ROWS FROM(
    UNNEST(@anchor_qs::SMALLINT[]),
    UNNEST(@anchor_rs::SMALLINT[])
) AS _params(
    anchor_q,
    anchor_r
)
WHERE
    _tile.city_id = @city_id
    AND _tile.anchor_q = _params.anchor_q
    AND _tile.anchor_r = _params.anchor_r
RETURNING
    _tile.q,
    _tile.r,
    _tile.anchor_q,
    _tile.anchor_r,
    _tile.artifact_id,
    _tile.artifact_status;

-- name: DestroyArtifact :many
UPDATE
    mutable.tiles AS _tile
SET
    artifact_id = @replacement_artifact_id,
    artifact_label = @replacement_artifact_label,
    artifact_status = @replacement_artifact_status,
    has_all_citizen_allocated = @has_all_citizen_allocated,
    anchor_q = _tile.q,
    anchor_r = _tile.r
WHERE
    _tile.city_id = @city_id
    AND _tile.anchor_q = @anchor_q
    AND _tile.anchor_r = @anchor_r
RETURNING
    _tile.super_kind_id,
    _tile.q,
    _tile.r,
    _tile.anchor_q,
    _tile.anchor_r,
    _tile.artifact_id,
    _tile.artifact_status;

-- name: DestroyPublicSystemsArtifact :many
UPDATE
    mutable.tiles AS _tile
SET
    super_kind_id = @replacement_super_kind_id,
    kind_id = @replacement_kind_id,
    artifact_id = @replacement_artifact_id,
    artifact_label = @replacement_artifact_label,
    artifact_status = @replacement_artifact_status,
    has_all_citizen_allocated = @has_all_citizen_allocated,
    anchor_q = _tile.q,
    anchor_r = _tile.r
WHERE
    _tile.city_id = @city_id
    AND _tile.anchor_q = @anchor_q
    AND _tile.anchor_r = @anchor_r
RETURNING
    _tile.super_kind_id,
    _tile.q,
    _tile.r,
    _tile.anchor_q,
    _tile.anchor_r,
    _tile.artifact_id,
    _tile.artifact_status;
