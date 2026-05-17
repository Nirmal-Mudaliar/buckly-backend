-- name: StoreRefreshToken :one
INSERT INTO refresh_tokens (
  user_id,
  token,
  expires_at,
  insert_ts
) VALUES (
  sqlc.arg(user_id),
  sqlc.arg(token),
  sqlc.arg(expires_at),
  sqlc.arg(insert_ts)
)
ON CONFLICT (user_id) 
DO UPDATE SET
  token = EXCLUDED.token,
  expires_at = EXCLUDED.expires_at,
  insert_ts = EXCLUDED.insert_ts
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT *
FROM refresh_tokens
WHERE token = sqlc.arg(token) AND expires_at > now()
LIMIT 1;