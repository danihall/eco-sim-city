-- name: UpsertKind :exec
INSERT INTO static.kinds(id, super_kind_id, name)
VALUES (@id, @super_kind_id, @name)
ON CONFLICT (id)
DO UPDATE SET name = EXCLUDED.name;
