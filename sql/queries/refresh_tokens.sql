-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: GetUserFromRefreshToken :one
SELECT u.id, u.created_at, u.updated_at, u.email, u.hashed_password, u.is_chirpy_red
FROM users u
JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.token = $1
    AND rt.expires_at > NOW()
    AND rt.revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;