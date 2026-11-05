-- name: UpsertSuperKind :exec
INSERT INTO static.super_kinds(id, name)
VALUES (@id, @name)
ON CONFLICT (id)
DO UPDATE SET name = EXCLUDED.name;
