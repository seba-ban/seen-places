CREATE EXTENSION IF NOT EXISTS postgis;

create table data_sources(
  filepath text unique primary key,
  type text not null,
  original_filename text not null,
  size bigint not null,
  created_at timestamp with time zone not null default current_timestamp
);

CREATE INDEX data_sources_type_idx
  ON data_sources (type);
CREATE INDEX data_sources_filepath_idx
  ON data_sources (filepath);

create table points(
  geom geography(Point, 4326) not null,
  visited timestamp with time zone not null,
  data_source_filepath text not null,
  CONSTRAINT fk_data_source
  FOREIGN KEY(data_source_filepath)
	REFERENCES data_sources(filepath)
);

CREATE INDEX points_geom_idx
  ON points
  USING GIST (geom);
CREATE INDEX points_visited_idx
  ON points (visited);

-- TODO: refactor
CREATE FUNCTION get_points(
    arg_visited_before TIMESTAMP WITH TIME ZONE = NULL, 
    arg_visited_after TIMESTAMP WITH TIME ZONE = NULL,
    arg_data_source_filepath TEXT = NULL,
    arg_data_source_type TEXT = NULL,
    arg_within_meters FLOAT = NULL,
    arg_within_meters_long_x FLOAT = NULL,
    arg_within_meters_lat_y FLOAT = NULL
) 
RETURNS TABLE(
	longitude_x FLOAT, 
	latitude_y FLOAT,
	filepath TEXT,
	type TEXT,
	geom GEOGRAPHY(Point, 4326)
) AS $$ 
	select
	    ST_X(p.geom::geometry)::float AS longitude_x,
	    ST_Y(p.geom::geometry)::float AS latitude_y,
		ds.filepath as filepath,
		ds.type as type,
		p.geom as geom
	from points p
	    join data_sources ds on p.data_source_filepath = ds.filepath
	where (
	        arg_visited_before is null
	        or p.visited <= arg_visited_before
	    )
	    and (
	        arg_visited_after is null
	        or p.visited > arg_visited_after
	    )
	    and (
	        arg_data_source_filepath is null
	        or ds.filepath = arg_data_source_filepath
	    )
	    and (
	        arg_data_source_type is null
	        or ds.type = arg_data_source_type
	    )
	    and ( (
	            arg_within_meters is null
	            or arg_within_meters_long_x is null
	            or arg_within_meters_lat_y is null
	        )
	        or ST_DWithin(
	            p.geom::geography,
	            ('SRID=4326;POINT(' || arg_within_meters_long_x || ' ' || arg_within_meters_lat_y || ')')::geography,
	            arg_within_meters
	        )
	    ) $$ LANGUAGE SQL;