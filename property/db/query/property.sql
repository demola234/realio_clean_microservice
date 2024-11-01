-- name: InsertProperty :one
INSERT INTO "Property" (
    id, title, description, price, type, address, zip_code, owner_id, images, 
    no_of_bed_rooms, no_of_bath_rooms, no_of_toilets, geo_location, status, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, 
    $10, $11, $12, $13, $14, now(), now()
)
RETURNING *;


-- name: GetPropertyByID :one
SELECT * FROM "Property"
WHERE id = $1;

-- name: ListProperties :many
SELECT * FROM "Property"
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateProperty :exec
UPDATE "Property"
SET 
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    price = COALESCE($4, price),
    type = COALESCE($5, type),
    address = COALESCE($6, address),
    zip_code = COALESCE($7, zip_code),
    images = COALESCE($8, images),
    no_of_bed_rooms = COALESCE($9, no_of_bed_rooms),
    no_of_bath_rooms = COALESCE($10, no_of_bath_rooms),
    no_of_toilets = COALESCE($11, no_of_toilets),
    geo_location = COALESCE($12, geo_location),
    status = COALESCE($13, status),
    updated_at = now()
WHERE id = $1;

-- name: DeleteProperty :exec
DELETE FROM "Property"
WHERE id = $1;


-- name: GetPropertiesByOwnerID :many
SELECT * FROM "Property"
WHERE owner_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;