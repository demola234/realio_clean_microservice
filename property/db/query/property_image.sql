-- name: ListPropertyImages :many
SELECT * FROM "PropertyImage"
WHERE property_id = $1
ORDER BY is_primary DESC, display_order ASC;

-- name: AddPropertyImage :one
INSERT INTO "PropertyImage" (
  id,
  property_id,
  url,
  caption,
  is_primary,
  display_order,
  room_type
) VALUES (
  gen_random_uuid(),
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdatePropertyImage :one
UPDATE "PropertyImage"
SET
  url = COALESCE($2, url),
  caption = COALESCE($3, caption),
  is_primary = COALESCE($4, is_primary),
  display_order = COALESCE($5, display_order),
  room_type = COALESCE($6, room_type),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePropertyImage :exec
DELETE FROM "PropertyImage"
WHERE id = $1;

-- name: SetPrimaryPropertyImage :exec
WITH update_primary AS (
  UPDATE "PropertyImage" pi1
  SET is_primary = false
  WHERE pi1.property_id = $1 AND pi1.id != $2
)
UPDATE "PropertyImage" pi2
SET is_primary = true
WHERE pi2.id = $2;