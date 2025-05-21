-- name: SaveSearchHistory :one
INSERT INTO "SearchHistory" (
  id,
  user_id,
  search_query,
  filters
) VALUES (
  gen_random_uuid(),
  $1, $2, $3
)
RETURNING *;

-- name: GetUserSearchHistory :many
SELECT * FROM "SearchHistory"
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: DeleteSearchHistory :exec
DELETE FROM "SearchHistory"
WHERE id = $1 AND user_id = $2;

-- name: ClearUserSearchHistory :exec
DELETE FROM "SearchHistory"
WHERE user_id = $1;