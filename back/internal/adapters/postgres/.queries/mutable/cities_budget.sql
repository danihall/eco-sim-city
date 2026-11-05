-- name: CreateCityBudget :copyfrom
INSERT INTO mutable.cities_budget (
    city_id,
    kind_id,
    tax_rate,
    funding_rate
)
VALUES (
    @city_id,
    @kind_id,
    @tax_rate,
    @funding_rate
);

-- name: GetCityKindBudget :one
SELECT
    tax_rate,
    funding_rate
FROM mutable.cities_budget
WHERE city_id = @city_id AND kind_id = @kind_id;

-- name: GetCityAllKindsBudget :many
SELECT
    kind_id,
    tax_rate,
    funding_rate
FROM mutable.cities_budget
WHERE city_id = @city_id;
