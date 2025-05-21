-- name: GetPropertyAmenities :many
SELECT
  a.*,
  pa.has_amenity,
  pa.notes
FROM "PropertyAmenity" pa
JOIN "Amenity" a ON pa.amenity_id = a.id
WHERE pa.property_id = $1
ORDER BY a.category, a.name;

-- name: AddPropertyAmenity :exec
INSERT INTO "PropertyAmenity" (
  property_id,
  amenity_id,
  has_amenity,
  notes
) VALUES (
  $1, $2, $3, $4
)
ON CONFLICT (property_id, amenity_id)
DO UPDATE
SET
  has_amenity = $3,
  notes = $4;

-- name: RemovePropertyAmenity :exec
DELETE FROM "PropertyAmenity"
WHERE property_id = $1 AND amenity_id = $2;

-- name: RemoveAllPropertyAmenities :exec
DELETE FROM "PropertyAmenity"
WHERE property_id = $1;