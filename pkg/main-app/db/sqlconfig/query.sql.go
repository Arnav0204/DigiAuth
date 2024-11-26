// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package db

import (
	"context"
)

const createConnection = `-- name: CreateConnection :exec
INSERT INTO connections (connection_id, id, my_mail_id, their_mail_id)
VALUES ($1, $2, $3, $4)
`

type CreateConnectionParams struct {
	ConnectionID string
	ID           int64
	MyMailID     string
	TheirMailID  string
}

func (q *Queries) CreateConnection(ctx context.Context, arg CreateConnectionParams) error {
	_, err := q.db.Exec(ctx, createConnection,
		arg.ConnectionID,
		arg.ID,
		arg.MyMailID,
		arg.TheirMailID,
	)
	return err
}

const createSchema = `-- name: CreateSchema :exec
INSERT INTO schemas (schema_id,credential_definition_id,schema_name,attributes)
VALUES ($1, $2, $3, $4)
`

type CreateSchemaParams struct {
	SchemaID               string
	CredentialDefinitionID string
	SchemaName             string
	Attributes             []string
}

func (q *Queries) CreateSchema(ctx context.Context, arg CreateSchemaParams) error {
	_, err := q.db.Exec(ctx, createSchema,
		arg.SchemaID,
		arg.CredentialDefinitionID,
		arg.SchemaName,
		arg.Attributes,
	)
	return err
}

const fetchConnections = `-- name: FetchConnections :many
SELECT connection_id, id, my_mail_id, their_mail_id
FROM connections
WHERE my_mail_id = $1
  AND their_mail_id = $2
`

type FetchConnectionsParams struct {
	MyMailID    string
	TheirMailID string
}

func (q *Queries) FetchConnections(ctx context.Context, arg FetchConnectionsParams) ([]Connection, error) {
	rows, err := q.db.Query(ctx, fetchConnections, arg.MyMailID, arg.TheirMailID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Connection
	for rows.Next() {
		var i Connection
		if err := rows.Scan(
			&i.ConnectionID,
			&i.ID,
			&i.MyMailID,
			&i.TheirMailID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getConnectionsByUserID = `-- name: GetConnectionsByUserID :many
SELECT connection_id, id, my_mail_id, their_mail_id 
FROM connections
WHERE id = $1
`

func (q *Queries) GetConnectionsByUserID(ctx context.Context, id int64) ([]Connection, error) {
	rows, err := q.db.Query(ctx, getConnectionsByUserID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Connection
	for rows.Next() {
		var i Connection
		if err := rows.Scan(
			&i.ConnectionID,
			&i.ID,
			&i.MyMailID,
			&i.TheirMailID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSchema = `-- name: GetSchema :many
SELECT schema_id, credential_definition_id, schema_name, attributes
FROM schemas
`

func (q *Queries) GetSchema(ctx context.Context) ([]Schema, error) {
	rows, err := q.db.Query(ctx, getSchema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Schema
	for rows.Next() {
		var i Schema
		if err := rows.Scan(
			&i.SchemaID,
			&i.CredentialDefinitionID,
			&i.SchemaName,
			&i.Attributes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSchemaById = `-- name: GetSchemaById :one
SELECT schema_id, credential_definition_id, schema_name, attributes
FROM schemas WHERE schema_id=$1
`

func (q *Queries) GetSchemaById(ctx context.Context, schemaID string) (Schema, error) {
	row := q.db.QueryRow(ctx, getSchemaById, schemaID)
	var i Schema
	err := row.Scan(
		&i.SchemaID,
		&i.CredentialDefinitionID,
		&i.SchemaName,
		&i.Attributes,
	)
	return i, err
}
