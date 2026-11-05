-- name: UpsertOrdinance :exec
INSERT INTO static.ordinances(id, name)
VALUES (@id, @name)
ON CONFLICT (id)
DO UPDATE SET name = EXCLUDED.name;
