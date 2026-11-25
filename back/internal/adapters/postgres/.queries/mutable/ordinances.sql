-- name: EnactOrdinances :exec
INSERT INTO
    mutable.ordinances_enacted (city_id, ordinance_id)
SELECT
    @city_id, UNNEST(@ordinance_ids::SMALLINT[])
ON CONFLICT DO NOTHING;

-- name: RevokeOrdinances :exec
DELETE FROM
    mutable.ordinances_enacted
WHERE
    city_id = @city_id;


-- name: GetEnactedOrdinances :many
SELECT
    ordinance_id
FROM
    mutable.ordinances_enacted
WHERE
    city_id = @city_id;
