-- name: ChirpyRedUpgrade :one

UPDATE users
SET is_chirpy_red = TRUE, updated_at = NOW()
where id = $1
returning id, created_at, updated_at, email, is_chirpy_red;