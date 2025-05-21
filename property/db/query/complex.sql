-- name: SearchPropertiesWithDetails :many
SELECT 
  p.*,
  pd.bedrooms,
  pd.bathrooms,
  pd.square_footage,
  pd.lot_size,
  pd.year_built,
  pd.stories,
  pd.garage_count,
  pd.has_basement,
  pd.has_attic,
  pn.walk_score,
  pn.transit_score,
  pn.school_rating,
  COALESCE(rev_stats.avg_rating, 0) as avg_rating,
  COALESCE(rev_stats.review_count, 0) as review_count
FROM "Property" p
LEFT JOIN "PropertyDetail" pd ON p.id = pd.property_id
LEFT JOIN "PropertyNeighborhood" pn ON p.id = pn.property_id
LEFT JOIN (
  SELECT 
    property_id, 
    AVG(overall_rating) as avg_rating,
    COUNT(*) as review_count
  FROM "Review"
  GROUP BY property_id
) rev_stats ON p.id = rev_stats.property_id
WHERE 
  ($1::property_category IS NULL OR p.category = $1) AND
  ($2::property_type IS NULL OR p.type = $2) AND
  ($3::property_status IS NULL OR p.status = $3) AND
  ($4::VARCHAR IS NULL OR p.city ILIKE '%' || $4 || '%') AND
  ($5::VARCHAR IS NULL OR p.state ILIKE '%' || $5 || '%') AND
  ($6::VARCHAR IS NULL OR p.country ILIKE '%' || $6 || '%') AND
  ($7::NUMERIC IS NULL OR p.price >= $7) AND
  ($8::NUMERIC IS NULL OR p.price <= $8) AND
  ($9::INT IS NULL OR pd.bedrooms >= $9) AND
  ($10::INT IS NULL OR pd.bathrooms >= $10) AND
  ($11::INT IS NULL OR pd.square_footage >= $11) AND
  ($12::INT IS NULL OR pd.year_built >= $12) AND
  ($13::INT IS NULL OR pd.garage_count >= $13) AND
  ($14::BOOLEAN IS NULL OR pd.has_basement = $14) AND
  ($15::BOOLEAN IS NULL OR pd.has_attic = $15) AND
  ($16::INT IS NULL OR pn.walk_score >= $16) AND
  ($17::INT IS NULL OR pn.school_rating >= $17)
ORDER BY 
  CASE WHEN $18 = 'price_asc' THEN p.price END ASC,
  CASE WHEN $18 = 'price_desc' THEN p.price END DESC,
  CASE WHEN $18 = 'newest' THEN p.created_at END DESC,
  CASE WHEN $18 = 'oldest' THEN p.created_at END ASC,
  CASE WHEN $18 = 'rating' THEN rev_stats.avg_rating END DESC,
  CASE WHEN $18 = 'most_reviewed' THEN rev_stats.review_count END DESC,
  p.created_at DESC
LIMIT $19
OFFSET $20;

-- name: GetPropertyWithAllDetails :one
SELECT 
  p.*,
  pd.*,
  pn.*,
  COALESCE(rev_stats.avg_rating, 0) as avg_rating,
  COALESCE(rev_stats.review_count, 0) as review_count
FROM "Property" p
LEFT JOIN "PropertyDetail" pd ON p.id = pd.property_id
LEFT JOIN "PropertyNeighborhood" pn ON p.id = pn.property_id
LEFT JOIN (
  SELECT 
    property_id, 
    AVG(overall_rating) as avg_rating,
    COUNT(*) as review_count
  FROM "Review"
  GROUP BY property_id
) rev_stats ON p.id = rev_stats.property_id
WHERE p.id = $1
LIMIT 1;

-- name: GetPropertiesByAmenity :many
SELECT 
  p.*
FROM "Property" p
JOIN "PropertyAmenity" pa ON p.id = pa.property_id
JOIN "Amenity" a ON pa.amenity_id = a.id
WHERE 
  pa.has_amenity = true AND
  a.id = $1 AND
  ($2::property_status IS NULL OR p.status = $2)
ORDER BY p.created_at DESC
LIMIT $3
OFFSET $4;

-- name: GetPropertiesByMultipleAmenities :many
SELECT 
  p.*,
  COUNT(pa.amenity_id) as matched_amenities
FROM "Property" p
JOIN "PropertyAmenity" pa ON p.id = pa.property_id
WHERE 
  pa.has_amenity = true AND
  pa.amenity_id = ANY($1::UUID[]) AND
  ($2::property_status IS NULL OR p.status = $2)
GROUP BY p.id
ORDER BY matched_amenities DESC, p.created_at DESC
LIMIT $3
OFFSET $4;