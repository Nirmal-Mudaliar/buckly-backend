-- name: CreateBucketListItem :one
INSERT INTO bucket_list_items (
  user_id,
  activity_tag_id,
  country_id,
  state_id,
  city_id,
  timeframe_start_date,
  timeframe_end_date,
  note,
  is_public,
  status,
  insert_ts,
  modified_ts
) VALUES (
  sqlc.arg(user_id),
  sqlc.arg(activity_tag_id),
  sqlc.arg(country_id),
  sqlc.arg(state_id),
  sqlc.arg(city_id),
  sqlc.arg(timeframe_start_date),
  sqlc.arg(timeframe_end_date),
  sqlc.arg(note),
  sqlc.arg(is_public),
  sqlc.arg(status),
  sqlc.arg(insert_ts),
  sqlc.arg(modified_ts)
) RETURNING *;

-- name: GetBucketListItemsByUserId :many
SELECT * FROM bucket_list_items
WHERE user_id = sqlc.arg(user_id)
ORDER BY modified_ts DESC;

-- name: GetBucketListItemById :one
SELECT * FROM bucket_list_items
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id)
LIMIT 1;

-- name: UpdateBucketListItem :one
UPDATE bucket_list_items
SET
  activity_tag_id = sqlc.arg(activity_tag_id),
  country_id = sqlc.arg(country_id),
  state_id = sqlc.arg(state_id),
  city_id = sqlc.arg(city_id),
  timeframe_start_date = sqlc.arg(timeframe_start_date),
  timeframe_end_date = sqlc.arg(timeframe_end_date),
  note = sqlc.arg(note),
  is_public = sqlc.arg(is_public),
  modified_ts = sqlc.arg(modified_ts)
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id)
RETURNING *;

-- name: DeleteBucketListItem :exec
DELETE FROM bucket_list_items
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);

-- name: FindMatchesForBucketListItem :many
SELECT * FROM bucket_list_items
WHERE activity_tag_id = sqlc.arg(activity_tag_id)
AND city_id = sqlc.arg(city_id)
AND timeframe_start_date <= sqlc.arg(timeframe_end_date)
AND timeframe_end_date >= sqlc.arg(timeframe_start_date)
AND is_public = true
AND user_id != sqlc.arg(user_id)
AND status = 'planned'
ORDER BY timeframe_start_date ASC;
