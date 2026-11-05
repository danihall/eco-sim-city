-- name: CreateCityHistoryFirstEntry :exec
INSERT INTO mutable.cities_state_history (city_id, ghg_1_level)
    VALUES (@city_id, @ghg_1_level);

-- name: GetAllCitiesLatestHistoryEntry :many
SELECT _history_entry.*
FROM mutable.cities _city
CROSS JOIN LATERAL (
    SELECT *
    FROM mutable.cities_state_history
    WHERE city_id = _city.id
    ORDER BY game_tick DESC
    LIMIT 1
) _history_entry;
