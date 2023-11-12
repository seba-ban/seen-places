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


-- name: GetGeoJson :one
with lines as (
        select 
                st_simplify(
                        st_makeline(geom::geometry order by visited asc),
                        0.01
                ) l
        from points p
        join data_sources ds on p.data_source_filepath = ds.filepath
        group by filepath
) select 
        st_asgeojson(
                st_union(l)   
        )::text geojson
from lines;
-- with lines as (
--         select 
--                 st_simplify(
--                         st_makeline(geom::geometry order by visited asc),
--                         0.01
--                 ) l
--         from points p
--         join data_sources ds on p.data_source_filepath = ds.filepath
--         group by filepath
-- ), lines_intersect as (
--         select
--                 unnest(ST_ClusterIntersecting(l::geometry)) li
--         from lines
-- ), lines_merged as (
--         select
--                 st_simplify(
--                         st_linemerge(li, false),
--                         0.01
--                 ) lm
--         from lines_intersect
-- )
-- select 
--         st_asgeojson(
--                 st_union(lm)   
--         )::text geojson
-- from lines_merged;
