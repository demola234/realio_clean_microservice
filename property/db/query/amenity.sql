-- name: ListAmenities :many
SELECT * FROM "Amenity"
WHERE ($1::amenity_category IS NULL OR category = $1)
ORDER BY name;

-- name: GetAmenityByID :one
SELECT * FROM "Amenity"
WHERE id = $1
LIMIT 1;

-- name: CreateAmenity :one
INSERT INTO "Amenity" (
  id,
  name,
  description,
  icon,
  category
) VALUES (
  gen_random_uuid(),
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateAmenity :one
UPDATE "Amenity"
SET
  name = COALESCE($2, name),
  description = COALESCE($3, description),
  icon = COALESCE($4, icon),
  category = COALESCE($5, category),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAmenity :exec
DELETE FROM "Amenity"
WHERE id = $1;