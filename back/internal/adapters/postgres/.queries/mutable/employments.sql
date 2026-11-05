-- name: DeleteJobsContracts :exec
DELETE FROM
    mutable.employments AS _contract
WHERE
    _contract.city_id = @city_id
AND
(
    (_contract.workers_anchor_q, _contract.workers_anchor_r) IN (
        SELECT worker_q, worker_r FROM ROWS FROM (
            UNNEST(@workers_anchor_qs::SMALLINT[]),
            UNNEST(@workers_anchor_rs::SMALLINT[])
        ) AS _row(worker_q, worker_r)
    )
    OR
    (_contract.employer_anchor_q, _contract.employer_anchor_r) IN (
        SELECT employer_q, employer_r FROM ROWS FROM (
            UNNEST(@employer_anchor_qs::SMALLINT[]),
            UNNEST(@employer_anchor_rs::SMALLINT[])
        ) AS _row(employer_q, employer_r)
    )
);

-- name: InsertJobsContracts :exec
INSERT INTO mutable.employments (
    city_id,
    worker_artifact_id,
    employer_anchor_q,
    employer_anchor_r,
    workers_anchor_q,
    workers_anchor_r,
    jobs_count
)
SELECT
    @city_id AS city_id,
    w_a AS worker_artifact_id,
    e_q AS employer_anchor_q,
    e_r AS employer_anchor_r,
    w_q AS workers_anchor_q,
    w_r AS workers_anchor_r,
    j_count AS jobs_count
FROM ROWS FROM (
    UNNEST(@worker_artifact_ids::SMALLINT[]),
    UNNEST(@employer_anchor_qs::SMALLINT[]),
    UNNEST(@employer_anchor_rs::SMALLINT[]),
    UNNEST(@workers_anchor_qs::SMALLINT[]),
    UNNEST(@workers_anchor_rs::SMALLINT[]),
    UNNEST(@jobs_count::SMALLINT[])
) AS _params(w_a, e_q, e_r, w_q, w_r, j_count)
ON CONFLICT (city_id, employer_anchor_q, employer_anchor_r, workers_anchor_q, workers_anchor_r)
DO UPDATE SET
    jobs_count = EXCLUDED.jobs_count;

-- name: SetAllCitizenAllocation :exec
UPDATE mutable.tiles AS _tile
SET
    has_all_citizen_allocated = _params.allocation
FROM ROWS FROM(
    UNNEST(@anchor_qs::SMALLINT[]),
    UNNEST(@anchor_rs::SMALLINT[]),
    UNNEST(@allocations::BOOLEAN[])
) AS _params(
    anchor_q,
    anchor_r,
    allocation
)
WHERE
    _tile.city_id = @city_id
    AND _tile.anchor_q = _params.anchor_q
    AND _tile.anchor_r = _params.anchor_r;

-- name: GetWorkerJobContractsFromChunk :many
SELECT
    id,
    worker_artifact_id,
    employer_anchor_q,
    employer_anchor_r,
    workers_anchor_q,
    workers_anchor_r,
    jobs_count
FROM mutable.employments
WHERE
    city_id = @city_id
    AND workers_anchor_q BETWEEN @q_top_left AND @q_bottom_right
    AND workers_anchor_r BETWEEN @r_top_left AND @r_bottom_right;
