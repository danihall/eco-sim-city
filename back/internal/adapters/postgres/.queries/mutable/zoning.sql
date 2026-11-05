-- name: UpdateTilesZones :many
UPDATE mutable.tiles AS _tile
SET
    super_kind_id = @new_super_kind_id,
    artifact_status = @new_artifact_status
FROM ROWS FROM(UNNEST(@qs::SMALLINT[]), UNNEST(@rs::SMALLINT[]))
    AS _params(q, r)
WHERE
    _tile.city_id = @city_id
    AND _tile.q = _params.q
    AND _tile.r = _params.r
    AND (
        _tile.super_kind_id = @natural_super_kind_id
        OR (
            _tile.super_kind_id != @public_systems_super_kind_id
            AND _tile.artifact_status = @old_artifact_status
        )
    )
RETURNING
    _tile.q,
    _tile.r,
    _tile.artifact_status;

-- name: RevertTilesZones :many
UPDATE
    mutable.tiles AS _tile
SET
    super_kind_id = @new_super_kind_id,
    artifact_status = @new_artifact_status
FROM ROWS FROM(UNNEST(@qs::SMALLINT[]), UNNEST(@rs::SMALLINT[]))
    AS _params(q, r)
WHERE
    _tile.city_id = @city_id
    AND _tile.q = _params.q
    AND _tile.r = _params.r
    AND _tile.super_kind_id != @natural_super_kind_id
    AND _tile.super_kind_id != @public_systems_super_kind_id
    AND _tile.artifact_status = @old_artifact_status
RETURNING
    _tile.q,
    _tile.r,
    _tile.artifact_status;
