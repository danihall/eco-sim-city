-- name: CreateCityBalance :exec
INSERT INTO mutable.cities_balances (city_id)
    VALUES (@city_id);

-- name: GetCityBalance :one
SELECT
    revenues,
    expenses
FROM mutable.cities_balances
WHERE city_id = @city_id;
