-- name: GetAllChirps :many
SELECT id, created_at, updated_at, body, user_id
FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpByChirpID :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: GetChirpsByUserID :many
SELECT *
FROM chirps
WHERE user_id = $1;