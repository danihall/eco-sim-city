-- name: CreateCityTilesSpatialEffects :copyfrom
INSERT INTO mutable.tiles_spatial_effects (
    city_id,
    q,
    r,
    ghg_1,
    crime_generic,
    crime_drug,
    crime_vss,
    land_value
)
VALUES (
    @city_id,
    @q,
    @r,
    @ghg_1,
    @crime_generic,
    @crime_drug,
    @crime_vss,
    @land_value
);

-- name: GetTilesSpatialEffects :many
SELECT
    _effects.ghg_1,
    _effects.crime_generic,
    _effects.crime_drug,
    _effects.crime_vss,
    _effects.land_value
FROM ROWS FROM(UNNEST(@qs::SMALLINT[]), UNNEST(@rs::SMALLINT[]))
    AS _params(q, r)
INNER JOIN mutable.tiles_spatial_effects _effects
ON _effects.city_id = @city_id
AND _effects.q = _params.q
AND _effects.r = _params.r;

-- name: UpdateTilesCriminality :exec
UPDATE mutable.tiles_spatial_effects AS _tile
SET
    crime_generic = _tile.crime_generic + _params.crime_generic_delta,
    crime_drug = _tile.crime_drug + _params.crime_drug_delta,
    crime_vss = _tile.crime_vss + _params.crime_vss_delta
FROM ROWS FROM(
    UNNEST(@qs::SMALLINT[]),
    UNNEST(@rs::SMALLINT[]),
    UNNEST(@crime_generic_deltas::SMALLINT[]),
    UNNEST(@crime_drug_deltas::SMALLINT[]),
    UNNEST(@crime_vss_deltas::SMALLINT[])
) AS _params(
    q,
    r,
    crime_generic_delta,
    crime_drug_delta,
    crime_vss_delta
)
WHERE
    _tile.city_id = @city_id
    AND _tile.q = _params.q
    AND _tile.r = _params.r;

-- name: UpdateTilesDesirability :exec
UPDATE mutable.tiles_spatial_effects AS _tile
SET
    land_value = _tile.land_value + _params.land_value_delta
FROM ROWS FROM(
    UNNEST(@qs::SMALLINT[]),
    UNNEST(@rs::SMALLINT[]),
    UNNEST(@land_value_deltas::SMALLINT[])
) AS _params(
    q,
    r,
    land_value_delta
)
WHERE
    _tile.city_id = @city_id
    AND _tile.q = _params.q
    AND _tile.r = _params.r;

-- name: UpdateTilesPollution :exec
UPDATE mutable.tiles_spatial_effects AS _tile
SET
    ghg_1 = _tile.ghg_1 + _params.ghg_1_delta
FROM ROWS FROM(
    UNNEST(@qs::SMALLINT[]),
    UNNEST(@rs::SMALLINT[]),
    UNNEST(@ghg_1_deltas::SMALLINT[])
) AS _params(
    q,
    r,
    ghg_1_delta
)
WHERE
    _tile.city_id = @city_id
    AND _tile.q = _params.q
    AND _tile.r = _params.r;
