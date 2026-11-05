-- name: GetUniqueArtifacts :many
SELECT
    artifact_id
FROM
    mutable.cities_flagships
WHERE
    city_id = @city_id;


-- name: AddUniqueArtifact :exec
INSERT INTO
    mutable.cities_flagships (city_id, artifact_id)
VALUES
    (@city_id, @artifact_id);

-- name: RemoveUniqueArtifact :exec
DELETE FROM
    mutable.cities_flagships
WHERE
    city_id = @city_id
    AND artifact_id = @artifact_id;
