-- name: CreateChirp :one
INSERT INTO chirps (body, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllChirps :many
SELECT * 
FROM CHIRPS
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT *
FROM CHIRPS
WHERE id = $1;

-- name: GetChirpsByAuthor :many
SELECT *
FROM CHIRPS
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;