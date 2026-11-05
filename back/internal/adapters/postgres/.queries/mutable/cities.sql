-- name: CreateCity :one
INSERT INTO mutable.cities (email, PASSWORD, name, mayor_name)
VALUES (@email, @password, @name, @mayor_name)
RETURNING id;

-- name: DeleteCity :one
DELETE FROM mutable.cities
WHERE id = @city_id
RETURNING id; -- needed to trigger an error if invalid id

-- name: CheckIfCityExists :one
SELECT EXISTS (
    SELECT 1
    FROM mutable.cities
    WHERE id = @city_id
) AS _exists;

-- name: GetCities :many
SELECT id
FROM mutable.cities;
