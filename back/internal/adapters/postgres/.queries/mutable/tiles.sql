-- name: CreateCityTiles :copyfrom
INSERT INTO mutable.tiles (
    city_id,
    super_kind_id,
    kind_id,
    artifact_id,
    artifact_label,
    artifact_status,
    q, r,
    anchor_r, anchor_q
)
VALUES (
    @city_id,
    @super_kind_id,
    @kind_id,
    @artifact_id,
    @artifact_label,
    @artifact_status,
    @q, @r,
    @anchor_r, @anchor_q
);

-- name: GetTilesAllPropsWithStats :many
SELECT
    _tile.super_kind_id,
    _tile.kind_id,
    _tile.artifact_id,
    _tile.artifact_label,
    _tile.artifact_status,
    _tile.q,
    _tile.r,
    _tile.anchor_r,
    _tile.anchor_q,
    _tile.has_all_citizen_allocated,
    _effects.ghg_1,
    _effects.crime_generic,
    _effects.crime_drug,
    _effects.crime_vss,
    _effects.land_value
FROM ROWS FROM(UNNEST(@qs::SMALLINT[]), UNNEST(@rs::SMALLINT[]))
    AS _params(q, r)
INNER JOIN mutable.tiles _tile
    ON _tile.city_id = @city_id
    AND _tile.q = _params.q
    AND _tile.r = _params.r
INNER JOIN mutable.tiles_spatial_effects _effects
ON _effects.city_id = @city_id
AND _effects.q = _params.q
AND _effects.r = _params.r
ORDER BY
    _tile.r ASC, _tile.q ASC
FOR UPDATE OF _tile;

-- name: GetTilesByAnchor :many
SELECT
    _tile.super_kind_id,
    _tile.kind_id,
    _tile.artifact_id,
    _tile.artifact_label,
    _tile.artifact_status,
    _tile.q,
    _tile.r,
    _tile.anchor_r,
    _tile.anchor_q
FROM
    mutable.tiles _tile
WHERE
    _tile.city_id = @city_id
    AND _tile.anchor_q = @q
    AND _tile.anchor_r = @r
ORDER BY
    _tile.r ASC, _tile.q ASC
FOR UPDATE OF _tile;

-- name: GetTilesFromCity :many
SELECT
    _tile.super_kind_id,
    _tile.artifact_id,
    _tile.artifact_status,
    _tile.q,
    _tile.r,
    _tile.anchor_r,
    _tile.anchor_q
FROM
    mutable.tiles _tile
WHERE
    _tile.city_id = @city_id
ORDER BY
    _tile.r ASC, _tile.q ASC;

-- name: GetTilesResidentialOrIndustrialWithEffects :many
SELECT
    _tile.super_kind_id,
    _tile.kind_id,
    _tile.artifact_id,
    _tile.q,
    _tile.r,
    _tile.anchor_q,
    _tile.anchor_r,
    _effects.ghg_1,
    _effects.crime_generic,
    _effects.crime_drug,
    _effects.crime_vss,
    _effects.land_value
FROM mutable.tiles _tile
INNER JOIN mutable.tiles_spatial_effects _effects
    ON _effects.city_id = _tile.city_id
    AND _effects.q = _tile.q
    AND _effects.r = _tile.r
WHERE
    _tile.city_id = @city_id
    AND _tile.artifact_status = @artifact_status
    AND _tile.q BETWEEN @q_top_left AND @q_bottom_right
    AND _tile.r BETWEEN @r_top_left AND @r_bottom_right
    AND _tile.super_kind_id = ANY(@super_kind_ids::SMALLINT[])
ORDER BY
    _tile.r ASC, _tile.q ASC
FOR UPDATE OF _tile;

-- name: GetAnchorTilesIndustrialOrPublicWithCitizenToAllocate :many
SELECT
    _tile.super_kind_id,
    _tile.artifact_id,
    _tile.q,
    _tile.r,
    _tile.anchor_r,
    _tile.anchor_q,
    COALESCE(SUM(_employer.jobs_count), 0)::SMALLINT AS jobs_count
FROM (
    SELECT *
    FROM mutable.tiles _t
    WHERE
        _t.city_id = @city_id
        AND _t.artifact_status = @artifact_status
        AND _t.q BETWEEN @q_top_left AND @q_bottom_right
        AND _t.r BETWEEN @r_top_left AND @r_bottom_right
        AND _t.anchor_q = _t.q
        AND _t.anchor_r = _t.r
        AND (
            _t.super_kind_id = @super_kind_industrial
            OR _t.super_kind_id = @super_kind_public_systems
        )
        AND _t.has_all_citizen_allocated = @at_capacity
    FOR UPDATE OF _t
) _tile
LEFT JOIN
    mutable.employments _employer
    ON _employer.city_id = _tile.city_id
    AND _employer.employer_anchor_q = _tile.anchor_q
    AND _employer.employer_anchor_r = _tile.anchor_r
GROUP BY
    _tile.super_kind_id,
    _tile.artifact_id,
    _tile.q,
    _tile.r,
    _tile.anchor_r,
    _tile.anchor_q
ORDER BY
    _tile.r ASC, _tile.q ASC;

-- name: GetAnchorTilesResidentialWithCitizenToAllocate :many
SELECT
    _tile.super_kind_id,
    _tile.artifact_id,
    _tile.q,
    _tile.r,
    _tile.anchor_r,
    _tile.anchor_q,
    COALESCE(SUM(_worker.jobs_count), 0)::SMALLINT AS jobs_count
FROM (
    SELECT *
    FROM mutable.tiles _t
    WHERE
        _t.city_id = @city_id
        AND _t.artifact_status = @artifact_status
        AND _t.q BETWEEN @q_top_left AND @q_bottom_right
        AND _t.r BETWEEN @r_top_left AND @r_bottom_right
        AND _t.anchor_q = _t.q
        AND _t.anchor_r = _t.r
        AND _t.super_kind_id = @super_kind_residential
        AND _t.has_all_citizen_allocated = @at_capacity
    FOR UPDATE OF _t
) _tile
LEFT JOIN
    mutable.employments _worker
    ON _worker.city_id = _tile.city_id
    AND _worker.workers_anchor_q = _tile.anchor_q
    AND _worker.workers_anchor_r = _tile.anchor_r
GROUP BY
    _tile.super_kind_id,
    _tile.artifact_id,
    _tile.q,
    _tile.r,
    _tile.anchor_r,
    _tile.anchor_q
ORDER BY
    _tile.r ASC, _tile.q ASC;

-- name: GetTransportNetworkTiles :many
SELECT
    q,
    r,
    artifact_id
FROM mutable.tiles
WHERE
    city_id = @city_id
    AND artifact_id = ANY(@artifact_ids::SMALLINT[]);
