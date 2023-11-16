-- name: CreateDataSource :one
INSERT INTO
    data_sources (
        filepath,
        type,
        original_filename,
        size
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetDataSourceByFilePath :one
SELECT * FROM data_sources WHERE filepath = $1;

-- name: GetDataSources :many
with dates as (
        select
            ds.filepath,
            min(p.visited)::timestamp start,
            max(p.visited)::timestamp end
from data_sources ds
    join points p on p.data_source_filepath = ds.filepath
group by ds.filepath
)
select
    ds.type,
    ds.filepath,
    ds.original_filename,
    ds.size,
    d.start,
    d.end
from data_sources ds
    join dates d on ds.filepath = d.filepath
where (
        sqlc.narg('type')::text is null
        or ds.type = sqlc.narg('type')::text
    )
    and (
        sqlc.narg('original_filename')::text is null
        or ds.original_filename = sqlc.narg('original_filename')::text
    )
    and (
        sqlc.narg('start_before')::timestamp is null
        or d.start < sqlc.narg('start_before')::timestamp
    )
    and (
        sqlc.narg('start_after')::timestamp is null
        or d.start > sqlc.narg('start_after')::timestamp
    )
order by d.start desc;

-- name: GetDataSourcesFromPolygon :many
with dates as (
        select
            ds.filepath,
            min(p.visited)::timestamp start,
            max(p.visited)::timestamp end
from data_sources ds
    join points p on p.data_source_filepath = ds.filepath
group by ds.filepath
)
select
    ds.type,
    ds.filepath,
    ds.original_filename,
    ds.size,
    d.start,
    d.end
from data_sources ds
    join dates d on ds.filepath = d.filepath
where ds.filepath in (
    select distinct data_source_filepath
    from points
    where st_within(
            geom::geometry, 
            st_polygon(
                    st_linefromtext(@linestring::text), 
                    4326
            )
    )
)
order by d.start desc;