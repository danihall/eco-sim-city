-- name: GetCityTraffic :many
SELECT
    q,
    r,
    SUM(foot_value)    AS foot_value,
    SUM(car_value)     AS car_value,
    SUM(bus_value)     AS bus_value,
    SUM(bicycle_value) AS bicycle_value,
    SUM(train_value)   AS train_value,
    SUM(subway_value)  AS subway_value
FROM
    mutable.traffic
WHERE
    city_id = @city_id
GROUP BY
    q, r;

-- name: DeleteCityTrafficFromChunk :exec
DELETE FROM
    mutable.traffic
WHERE
    city_id = @city_id
    AND workers_anchor_q BETWEEN @q_top_left AND @q_bottom_right
    AND workers_anchor_r BETWEEN @r_top_left AND @r_bottom_right;




-- name: InsertTrafficBatch :copyfrom
INSERT INTO mutable.traffic (
    city_id,
    employment_id,
    q,
    r,
    workers_anchor_q,
    workers_anchor_r,
    foot_value,
    car_value,
    bus_value,
    bicycle_value,
    train_value,
    subway_value
)
VALUES (
    @city_id,
    @employment_id,
    @q,
    @r,
    @workers_anchor_q,
    @workers_anchor_r,
    @foot_value,
    @car_value,
    @bus_value,
    @bicycle_value,
    @train_value,
    @subway_value
);
