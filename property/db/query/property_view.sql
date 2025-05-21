-- name: RecordPropertyView :one
INSERT INTO "PropertyView" (
  id,
  property_id,
  user_id
) VALUES (
  gen_random_uuid(),
  $1, $2
)
RETURNING *;

-- name: GetPropertyViewStats :one
SELECT
  COUNT(*) as total_views,
  COUNT(DISTINCT user_id) as unique_views
FROM "PropertyView"
WHERE property_id = $1;

-- name: GetRecentlyViewedProperties :many
SELECT DISTINCT ON (pv.property_id)
  p.*,
  pv.viewed_at
FROM "PropertyView" pv
JOIN "Property" p ON pv.property_id = p.id
WHERE pv.user_id = $1
ORDER BY pv.property_id, pv.viewed_at DESC
LIMIT $2;