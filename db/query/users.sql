-- name: CreateUser :one
INSERT INTO users (
  email,
  hashed_password,
  password_changed_at,  
  first_name,
  last_name,
  profile_image_path,
  salt,
  last_login,
  created_at,
  created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUserByEmail :one
UPDATE users
SET
  email = COALESCE(sqlc.narg(email), email),
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  profile_image_path = COALESCE(sqlc.narg(profile_image_path), profile_image_path),
  salt = COALESCE(sqlc.narg(salt), salt),
  last_login = COALESCE(sqlc.narg(last_login), last_login),
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  created_at = COALESCE(sqlc.narg(created_at), created_at),
  created_by = COALESCE(sqlc.narg(created_by), created_by)
WHERE
  email = sqlc.arg(email)
RETURNING *;