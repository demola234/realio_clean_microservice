-- name: GetPropertyAvailability :many
SELECT * FROM "PropertyAvailability"
WHERE property_id = $1
  AND date BETWEEN $2 AND $3
ORDER BY date;

-- name: SetPropertyAvailability :one
INSERT INTO "PropertyAvailability" (
  id,
  property_id,
  date,
  is_available,
  price_override,
  min_nights,
  max_nights
) VALUES (
  gen_random_uuid(),
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT (property_id, date)
DO UPDATE
SET
  is_available = $3,
  price_override = $4,
  min_nights = $5,
  max_nights = $6,
  updated_at = now()
RETURNING *;

-- name: BulkSetPropertyAvailabilityStatus :exec
WITH dates AS (
  SELECT generate_series($2::date, $3::date, '1 day'::interval)::date as date
)
INSERT INTO "PropertyAvailability" (
  id,
  property_id,
  date,
  is_available
)
SELECT
  gen_random_uuid(),
  $1,
  dates.date,
  $4
FROM dates
ON CONFLICT (property_id, date)
DO UPDATE
SET
  is_available = $4,
  updated_at = now();