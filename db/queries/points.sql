-- name: CreatePoints :copyfrom
INSERT INTO points (geom, visited, data_source_filepath) VALUES ($1, $2, $3);

-- name: GetDataSourcePointsCount :one
select count(*)
from get_points(arg_data_source_filepath := sqlc.arg('data_source_filepath'));

-- name: GetLineStrings :many
select
    st_asgeojson(
        st_makeline(geom::geometry)
    )::text lines
from get_points(
    sqlc.narg('visited_before'),
    sqlc.narg('visited_after'),
    sqlc.narg('data_source_filepath'),
    sqlc.narg('data_source_type'),
    sqlc.narg('within_meters'),
    sqlc.narg('within_meters_long_x'),
    sqlc.narg('within_meters_lat_y')
)
group by filepath::text;

-- name: GetLineByFilepath :one
select
    st_asgeojson(
        st_makeline(geom::geometry order by visited desc)
    )::text
from points
where data_source_filepath = $1;

-- name: GetGeoJson :one
with a as (
        select st_simplify(st_makeline(geom::geometry order by visited), 0.001) c
        from points
        group by data_source_filepath
) select st_asgeojson(st_linemerge(st_union(c)))::text from a;
