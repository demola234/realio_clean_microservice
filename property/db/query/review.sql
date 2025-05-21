-- name: GetPropertyReviews :many
SELECT * FROM "Review"
WHERE property_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetPropertyReviewStats :one
SELECT
  COUNT(*) as total_reviews,
  AVG(overall_rating) as avg_overall,
  AVG(location_rating) as avg_location,
  AVG(value_rating) as avg_value,
  AVG(accuracy_rating) as avg_accuracy,
  AVG(communication_rating) as avg_communication,
  AVG(cleanliness_rating) as avg_cleanliness,
  AVG(check_in_rating) as avg_check_in
FROM "Review"
WHERE property_id = $1;

-- name: CreateReview :one
INSERT INTO "Review" (
  id,
  booking_id,
  property_id,
  reviewer_id,
  reviewed_id,
  overall_rating,
  location_rating,
  value_rating,
  accuracy_rating,
  communication_rating,
  cleanliness_rating,
  check_in_rating,
  comment,
  type
) VALUES (
  gen_random_uuid(),
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM "Review"
WHERE id = $1;