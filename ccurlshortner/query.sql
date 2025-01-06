-- name: CreateUrl :one
INSERT INTO
    url (key, longurl, shorturl)
VALUES
    (?, ?, ?) RETURNING *;

-- name: SelectShortUrl :one
SELECT
    longurl
FROM
    url
WHERE
    shorturl = ?;

-- name: DeleteUrl :exec
DELETE FROM url
WHERE
    key = ?;
