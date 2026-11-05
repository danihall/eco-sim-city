-- name: EnactOrdinance :exec
INSERT INTO
    mutable.ordinances_enacted (city_id, ordinance_id)
VALUES (@city_id, @ordinance_id);

-- name: RevokeOrdinance :exec
DELETE FROM
    mutable.ordinances_enacted
WHERE
    city_id = @city_id
    AND ordinance_id = @ordinance_id;

-- name: GetEnactedOrdinances :many
SELECT
    ordinance_id
FROM
    mutable.ordinances_enacted
WHERE
    city_id = @city_id;
