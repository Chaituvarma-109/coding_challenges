// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package urldb

import (
	"context"
)

const createUrl = `-- name: CreateUrl :one
INSERT INTO
    url (key, longurl, shorturl)
VALUES
    (?, ?, ?) RETURNING "key", longurl, shorturl
`

type CreateUrlParams struct {
	Key      string `json:"key"`
	Longurl  string `json:"longurl"`
	Shorturl string `json:"shorturl"`
}

func (q *Queries) CreateUrl(ctx context.Context, arg CreateUrlParams) (Url, error) {
	row := q.db.QueryRowContext(ctx, createUrl, arg.Key, arg.Longurl, arg.Shorturl)
	var i Url
	err := row.Scan(&i.Key, &i.Longurl, &i.Shorturl)
	return i, err
}

const deleteUrl = `-- name: DeleteUrl :exec
DELETE FROM url
WHERE
    key = ?
`

func (q *Queries) DeleteUrl(ctx context.Context, key string) error {
	_, err := q.db.ExecContext(ctx, deleteUrl, key)
	return err
}

const selectShortUrl = `-- name: SelectShortUrl :one
SELECT
    longurl
FROM
    url
WHERE
    shorturl = ?
`

func (q *Queries) SelectShortUrl(ctx context.Context, shorturl string) (string, error) {
	row := q.db.QueryRowContext(ctx, selectShortUrl, shorturl)
	var longurl string
	err := row.Scan(&longurl)
	return longurl, err
}