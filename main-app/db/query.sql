-- name: GetConnectionsByUserID :many
SELECT * 
FROM connections
WHERE id = $1;

-- name: CreateConnection :exec
INSERT INTO connections (connection_id, id, alias, my_role)
VALUES ($1, $2, $3, $4);