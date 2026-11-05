-- name: UpsertArtifact :exec
INSERT INTO static.artifacts(id, super_kind_id, kind_id, name)
VALUES (@id, @super_kind_id, @kind_id, @name)
ON CONFLICT (id)
DO UPDATE SET name = EXCLUDED.name;
