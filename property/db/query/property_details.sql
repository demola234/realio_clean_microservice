-- name: GetPropertyDetailByPropertyID :one
SELECT * FROM "PropertyDetail"
WHERE property_id = $1
LIMIT 1;

-- name: CreatePropertyDetail :one
INSERT INTO "PropertyDetail" (
  id,
  property_id,
  bedrooms,
  bathrooms,
  toilets,
  square_footage,
  lot_size,
  year_built,
  stories,
  garage_count,
  has_basement,
  has_attic,
  heating_system,
  cooling_system,
  water_source,
  sewer_type,
  roof_type,
  exterior_material,
  foundation_type,
  pool_type,
  geo_location
) VALUES (
  gen_random_uuid(),
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
  $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
)
RETURNING *;

-- name: UpdatePropertyDetail :one
UPDATE "PropertyDetail"
SET
  bedrooms = COALESCE($2, bedrooms),
  bathrooms = COALESCE($3, bathrooms),
  toilets = COALESCE($4, toilets),
  square_footage = COALESCE($5, square_footage),
  lot_size = COALESCE($6, lot_size),
  year_built = COALESCE($7, year_built),
  stories = COALESCE($8, stories),
  garage_count = COALESCE($9, garage_count),
  has_basement = COALESCE($10, has_basement),
  has_attic = COALESCE($11, has_attic),
  heating_system = COALESCE($12, heating_system),
  cooling_system = COALESCE($13, cooling_system),
  water_source = COALESCE($14, water_source),
  sewer_type = COALESCE($15, sewer_type),
  roof_type = COALESCE($16, roof_type),
  exterior_material = COALESCE($17, exterior_material),
  foundation_type = COALESCE($18, foundation_type),
  pool_type = COALESCE($19, pool_type),
  geo_location = COALESCE($20, geo_location),
  updated_at = now()
WHERE property_id = $1
RETURNING *;