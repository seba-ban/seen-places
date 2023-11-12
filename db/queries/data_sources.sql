-- name: CreateDataSource :one
INSERT INTO data_sources (
  filepath, type, original_filename, size
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetDataSourceByFilePath :one
SELECT * FROM data_sources
WHERE filepath = $1;
