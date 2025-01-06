-- name: CreateUrl :one
INSERT INTO
    url (key, longurl, shorturl)
VALUES
    (?, ?, ?) RETURNING *;

-- name: SelectLongUrl :one
SELECT
    longurl
FROM
    url
WHERE
    key = ?;

-- name: DeleteUrl :exec
DELETE FROM url
WHERE
    key = ?;
