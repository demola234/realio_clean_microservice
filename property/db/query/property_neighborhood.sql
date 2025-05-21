-- name: GetPropertyNeighborhood :one
SELECT * FROM "PropertyNeighborhood"
WHERE property_id = $1
LIMIT 1;

-- name: CreatePropertyNeighborhood :one
INSERT INTO "PropertyNeighborhood" (
  id,
  property_id,
  school_district,
  school_rating,
  crime_rate,
  walk_score,
  transit_score,
  bike_score,
  nearby_locations
) VALUES (
  gen_random_uuid(),
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdatePropertyNeighborhood :one
UPDATE "PropertyNeighborhood"
SET
  school_district = COALESCE($2, school_district),
  school_rating = COALESCE($3, school_rating),
  crime_rate = COALESCE($4, crime_rate),
  walk_score = COALESCE($5, walk_score),
  transit_score = COALESCE($6, transit_score),
  bike_score = COALESCE($7, bike_score),
  nearby_locations = COALESCE($8, nearby_locations),
  updated_at = now()
WHERE property_id = $1
RETURNING *;