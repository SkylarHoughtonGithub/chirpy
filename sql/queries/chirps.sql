-- name: CreateChirp :one
INSERT INTO chirps (body, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllChirp :one
SELECT * 
FROM CHIRPS
ORDER BY created_at ASC;