-- name: AddRefreshToken :one

INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES (
    $1, -- token
    NOW(), -- created_at
    NOW(), -- updated_at
    NOW() + INTERVAL '60 days', -- expires_at (set to 60 days from now)
    NULL, -- revoked_at (initially null)
    $2  -- user_id
)
RETURNING *;