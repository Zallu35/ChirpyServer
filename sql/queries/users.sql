-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: PostChirp :one
INSERT INTO posts (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: RetrieveChirps :many
SELECT * FROM posts 
ORDER BY created_at ASC;

-- name: RetrieveChirpsByAuthor :many
SELECT * FROM posts WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetSingleChirp :one
Select * FROM posts WHERE id = $1;

-- name: LoginQuery :one
SELECT * FROM users WHERE email = $1;

-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT user_id, expires_at, revoked_at FROM refresh_tokens WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW() WHERE token = $1;

-- name: UpdateCredentials :one
UPDATE users SET email = $2, hashed_password = $3, updated_at = NOW() WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red;

-- name: DeleteChirp :exec
DELETE FROM posts WHERE id = $1;

-- name: UpgradeChirpyRed :exec
UPDATE users SET is_chirpy_red = true WHERE id = $1;