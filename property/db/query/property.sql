
-- name: GetPropertyByID :one
SELECT * FROM "Property"
WHERE id = $1
LIMIT 1;

-- name: ListProperties :many
SELECT 
  p.id, p.title, p.description, p.price, p.category, p.type, p.address, p.city, p.state, p.country, p.zip_code, p.owner_id, p.status, p.created_at, p.updated_at,
  pd.bedrooms,
  pd.bathrooms,
  pd.square_footage
FROM "Property" p
LEFT JOIN "PropertyDetail" pd ON p.id = pd.property_id
WHERE 
  (sqlc.narg('category')::property_category IS NULL OR p.category = sqlc.narg('category')) AND
  (sqlc.narg('type')::property_type IS NULL OR p.type = sqlc.narg('type')) AND
  (sqlc.narg('status')::property_status IS NULL OR p.status = sqlc.narg('status')) AND
  (sqlc.narg('city')::VARCHAR IS NULL OR p.city = sqlc.narg('city')) AND
  (sqlc.narg('state')::VARCHAR IS NULL OR p.state = sqlc.narg('state')) AND
  (sqlc.narg('country')::VARCHAR IS NULL OR p.country = sqlc.narg('country')) AND
  (sqlc.narg('min_price')::NUMERIC IS NULL OR p.price >= sqlc.narg('min_price')) AND
  (sqlc.narg('max_price')::NUMERIC IS NULL OR p.price <= sqlc.narg('max_price')) AND
  (sqlc.narg('min_bedrooms')::INT IS NULL OR pd.bedrooms >= sqlc.narg('min_bedrooms')) AND
  (sqlc.narg('min_bathrooms')::INT IS NULL OR pd.bathrooms >= sqlc.narg('min_bathrooms'))
ORDER BY p.created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateProperty :one
INSERT INTO "Property" (
  id,
  title,
  description,
  price,
  category,
  type,
  address,
  city,
  state,
  country,
  zip_code,
  owner_id,
  status
) VALUES (
  gen_random_uuid(),
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: UpdateProperty :one
UPDATE "Property"
SET
  title = COALESCE($2, title),
  description = COALESCE($3, description),
  price = COALESCE($4, price),
  category = COALESCE($5, category),
  type = COALESCE($6, type),
  address = COALESCE($7, address),
  city = COALESCE($8, city),
  state = COALESCE($9, state),
  country = COALESCE($10, country),
  zip_code = COALESCE($11, zip_code),
  status = COALESCE($12, status),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdatePropertyStatus :one
UPDATE "Property"
SET
  status = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProperty :exec
DELETE FROM "Property"
WHERE id = $1;

-- name: CountProperties :one
SELECT COUNT(*) FROM "Property"
WHERE 
  ($1::property_category IS NULL OR category = $1) AND
  ($2::property_type IS NULL OR type = $2) AND
  ($3::property_status IS NULL OR status = $3) AND
  ($4::VARCHAR IS NULL OR city = $4) AND
  ($5::VARCHAR IS NULL OR state = $5) AND
  ($6::VARCHAR IS NULL OR country = $6) AND
  ($7::NUMERIC IS NULL OR price >= $7) AND
  ($8::NUMERIC IS NULL OR price <= $8);

-- name: GetPropertiesByOwner :many
SELECT * FROM "Property"
WHERE owner_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;