-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  phone_no,
  date_of_birth,
  gender,
  password_hash,
  is_phone_verified,
  insert_ts,
  modified_ts 
) VALUES (
  sqlc.arg(first_name),
  sqlc.arg(last_name),  
  sqlc.arg(email),
  sqlc.arg(phone_no),
  sqlc.arg(date_of_birth),
  sqlc.arg(gender),
  sqlc.arg(password_hash),
  sqlc.arg(is_phone_verified),
  sqlc.arg(insert_ts),
  sqlc.arg(modified_ts)
) RETURNING *;