-- name: GetCurrentCityState :one
SELECT
    game_tick,
    electricity_consumption,
    water_consumption,
    ghg_1_level,
    ghg_2_level,
    ghg_3_level,
    climate_change_level,
    children_count,
    adults_count,
    elderly_count,
    citizens_count,
    children_education_level,
    adults_education_level,
    elderly_education_level,
    citizens_education_level,
    children_health_level,
    adults_health_level,
    elderly_health_level,
    citizens_health_level
FROM computed.current_cities_state
WHERE city_id = @city_id;

-- name: PopulateCitiesCurrentState :copyfrom
INSERT INTO computed.current_cities_state (
    city_id,
    game_tick,
    electricity_consumption,
    water_consumption,
    ghg_1_level,
    ghg_2_level,
    ghg_3_level,
    children_count,
    adults_count,
    elderly_count,
    children_education_level,
    adults_education_level,
    elderly_education_level,
    children_health_level,
    adults_health_level,
    elderly_health_level
)
VALUES (
    @city_id,
    @game_tick,
    @electricity_consumption,
    @water_consumption,
    @ghg_1_level,
    @ghg_2_level,
    @ghg_3_level,
    @children_count,
    @adults_count,
    @elderly_count,
    @children_education_level,
    @adults_education_level,
    @elderly_education_level,
    @children_health_level,
    @adults_health_level,
    @elderly_health_level
);

-- name: ClearCitiesCurrentState :exec
TRUNCATE TABLE computed.current_cities_state;
