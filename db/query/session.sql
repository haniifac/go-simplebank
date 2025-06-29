-- name: CreateSession :one
INSERT INTO sessions (
    id,
    username,
    refresh_token,
    user_agent,
    ip_address,
    created_at,
    expires_at,
    is_blocked
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions WHERE id = $1 LIMIT 1;