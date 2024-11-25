-- name: GetConnectionsByUserID :many
SELECT * 
FROM connections
WHERE id = $1;

-- name: CreateConnection :exec
INSERT INTO connections (connection_id, id, alias, my_role)
VALUES ($1, $2, $3, $4);

-- name: CreateSchema :exec
INSERT INTO schemas (schema_id,credential_definition_id,schema_name,attributes)
VALUES ($1, $2, $3, $4);

-- name: GetSchema :many
SELECT *
FROM schemas;

-- name: GetSchemaById :one
SELECT *
FROM schemas WHERE schema_id=$1;