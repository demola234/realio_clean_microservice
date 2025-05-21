CREATE TYPE property_category AS ENUM ('Rent', 'Sale', 'Buy', 'Lease', 'Land', 'Mortgage');
CREATE TYPE property_type AS ENUM ('House', 'Apartment', 'Land', 'Manufactured', 'Townhome', 'Multi-family', 'Condo', 'Co-op', 'Lot');
CREATE TYPE property_status AS ENUM ('Available', 'Sold', 'Rented', 'Pending');
CREATE TYPE amenity_category AS ENUM ('Feature', 'Amenity', 'Appliance', 'Utility', 'Parking');
CREATE TYPE review_type AS ENUM ('Property', 'Host', 'Guest');


CREATE TABLE "Property" (
  "id" UUID PRIMARY KEY,
  "title" VARCHAR NOT NULL,
  "description" TEXT,
  "price" NUMERIC NOT NULL,
  "category" property_category NOT NULL,
  "type" property_type NOT NULL,
  "address" VARCHAR NOT NULL,
  "city" VARCHAR NOT NULL,
  "state" VARCHAR NOT NULL,
  "country" VARCHAR NOT NULL,
  "zip_code" VARCHAR,
  "owner_id" UUID,
  "status" property_status NOT NULL DEFAULT 'Available',
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

COMMENT ON COLUMN "Property"."id" IS 'Primary key';
COMMENT ON COLUMN "Property"."title" IS 'Property title';
COMMENT ON COLUMN "Property"."description" IS 'Detailed description';
COMMENT ON COLUMN "Property"."price" IS 'Price of the property';
COMMENT ON COLUMN "Property"."category" IS 'Rent, Sale, Buy, Lease, Land, Mortgage, etc.';
COMMENT ON COLUMN "Property"."type" IS 'House, Apartment, Land, etc.';
COMMENT ON COLUMN "Property"."address" IS 'Address details';
COMMENT ON COLUMN "Property"."city" IS 'City';
COMMENT ON COLUMN "Property"."state" IS 'State/Province';
COMMENT ON COLUMN "Property"."country" IS 'Country';
COMMENT ON COLUMN "Property"."zip_code" IS 'Zip code/Postal code';
COMMENT ON COLUMN "Property"."owner_id" IS 'Reference to the user (seller), external microservice';
COMMENT ON COLUMN "Property"."status" IS 'Available, Sold, Rented, etc.';

CREATE TABLE "PropertyDetail" (
  "id" UUID PRIMARY KEY,
  "property_id" UUID NOT NULL,
  "bedrooms" INT,
  "bathrooms" INT,
  "toilets" INT,
  "square_footage" NUMERIC,
  "lot_size" NUMERIC,
  "year_built" INT,
  "stories" INT,
  "garage_count" INT,
  "has_basement" BOOLEAN DEFAULT false,
  "has_attic" BOOLEAN DEFAULT false,
  "heating_system" VARCHAR,
  "cooling_system" VARCHAR,
  "water_source" VARCHAR,
  "sewer_type" VARCHAR,
  "roof_type" VARCHAR,
  "exterior_material" VARCHAR,
  "foundation_type" VARCHAR,
  "pool_type" VARCHAR,
  "geo_location" JSON,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  FOREIGN KEY ("property_id") REFERENCES "Property"("id") ON DELETE CASCADE
);

COMMENT ON TABLE "PropertyDetail" IS 'Detailed information about property features';
COMMENT ON COLUMN "PropertyDetail"."bedrooms" IS 'Number of bedrooms';
COMMENT ON COLUMN "PropertyDetail"."bathrooms" IS 'Number of bathrooms';
COMMENT ON COLUMN "PropertyDetail"."toilets" IS 'Number of toilets';
COMMENT ON COLUMN "PropertyDetail"."square_footage" IS 'Total living area in square feet/meters';
COMMENT ON COLUMN "PropertyDetail"."lot_size" IS 'Size of the lot in square feet/meters';
COMMENT ON COLUMN "PropertyDetail"."year_built" IS 'Year the property was built';
COMMENT ON COLUMN "PropertyDetail"."stories" IS 'Number of stories/floors';
COMMENT ON COLUMN "PropertyDetail"."garage_count" IS 'Number of garage spaces';
COMMENT ON COLUMN "PropertyDetail"."geo_location" IS 'Latitude & longitude (Optional)';

CREATE TABLE "Amenity" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "description" TEXT,
  "icon" VARCHAR,
  "category" amenity_category NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

COMMENT ON TABLE "Amenity" IS 'Master list of all possible amenities, features, appliances, utilities, and parking options';
COMMENT ON COLUMN "Amenity"."name" IS 'Name of the amenity';
COMMENT ON COLUMN "Amenity"."description" IS 'Detailed description of the amenity';
COMMENT ON COLUMN "Amenity"."icon" IS 'Icon reference for UI display';
COMMENT ON COLUMN "Amenity"."category" IS 'Type of amenity (Feature, Amenity, Appliance, Utility, Parking)';

CREATE TABLE "PropertyAmenity" (
  "property_id" UUID NOT NULL,
  "amenity_id" UUID NOT NULL,
  "has_amenity" BOOLEAN DEFAULT true,
  "notes" TEXT,
  PRIMARY KEY ("property_id", "amenity_id"),
  FOREIGN KEY ("property_id") REFERENCES "Property"("id") ON DELETE CASCADE,
  FOREIGN KEY ("amenity_id") REFERENCES "Amenity"("id") ON DELETE CASCADE
);

COMMENT ON TABLE "PropertyAmenity" IS 'Links properties to their amenities with additional details';
COMMENT ON COLUMN "PropertyAmenity"."has_amenity" IS 'Indicates if property has this amenity (allows for explicitly noting absence)';
COMMENT ON COLUMN "PropertyAmenity"."notes" IS 'Additional notes about this amenity for this property';

CREATE TABLE "PropertyImage" (
  "id" UUID PRIMARY KEY,
  "property_id" UUID NOT NULL,
  "url" VARCHAR NOT NULL,
  "caption" VARCHAR,
  "is_primary" BOOLEAN DEFAULT false,
  "display_order" INT DEFAULT 0,
  "room_type" VARCHAR,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  FOREIGN KEY ("property_id") REFERENCES "Property"("id") ON DELETE CASCADE
);

COMMENT ON TABLE "PropertyImage" IS 'Property images with metadata';
COMMENT ON COLUMN "PropertyImage"."caption" IS 'Optional caption for the image';
COMMENT ON COLUMN "PropertyImage"."is_primary" IS 'Whether this is the main/featured image';
COMMENT ON COLUMN "PropertyImage"."display_order" IS 'Order to display images in';
COMMENT ON COLUMN "PropertyImage"."room_type" IS 'Type of room shown (e.g., Kitchen, Bathroom, Exterior)';

CREATE TABLE "PropertyAvailability" (
  "id" UUID PRIMARY KEY,
  "property_id" UUID NOT NULL,
  "date" DATE NOT NULL,
  "is_available" BOOLEAN DEFAULT true,
  "price_override" NUMERIC,
  "min_nights" INT,
  "max_nights" INT,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  FOREIGN KEY ("property_id") REFERENCES "Property"("id") ON DELETE CASCADE,
  UNIQUE ("property_id", "date")
);

COMMENT ON TABLE "PropertyAvailability" IS 'Calendar availability for rental properties';
COMMENT ON COLUMN "PropertyAvailability"."price_override" IS 'Special pricing for specific dates';
COMMENT ON COLUMN "PropertyAvailability"."min_nights" IS 'Minimum nights for bookings that include this date';
COMMENT ON COLUMN "PropertyAvailability"."max_nights" IS 'Maximum nights for bookings that include this date';

CREATE TABLE "Review" (
  "id" UUID PRIMARY KEY,
  "booking_id" UUID,
  "property_id" UUID NOT NULL,
  "reviewer_id" UUID NOT NULL,
  "reviewed_id" UUID,
  "overall_rating" NUMERIC(2,1) NOT NULL CHECK (overall_rating >= 1 AND overall_rating <= 5),
  "location_rating" NUMERIC(2,1) CHECK (location_rating >= 1 AND location_rating <= 5),
  "value_rating" NUMERIC(2,1) CHECK (value_rating >= 1 AND value_rating <= 5),
  "accuracy_rating" NUMERIC(2,1) CHECK (accuracy_rating >= 1 AND accuracy_rating <= 5),
  "communication_rating" NUMERIC(2,1) CHECK (communication_rating >= 1 AND communication_rating <= 5),
  "cleanliness_rating" NUMERIC(2,1) CHECK (cleanliness_rating >= 1 AND cleanliness_rating <= 5),
  "check_in_rating" NUMERIC(2,1) CHECK (check_in_rating >= 1 AND check_in_rating <= 5),
  "comment" TEXT,
  "type" review_type NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  FOREIGN KEY ("property_id") REFERENCES "Property"("id") ON DELETE CASCADE
);

COMMENT ON TABLE "Review" IS 'Reviews for properties with detailed ratings';
COMMENT ON COLUMN "Review"."booking_id" IS 'External reference to booking from the booking microservice';
COMMENT ON COLUMN "Review"."reviewer_id" IS 'User who wrote the review (from auth microservice)';
COMMENT ON COLUMN "Review"."reviewed_id" IS 'User being reviewed, if applicable (host/guest)';
COMMENT ON COLUMN "Review"."overall_rating" IS 'Overall rating score';
COMMENT ON COLUMN "Review"."location_rating" IS 'Rating for location';
COMMENT ON COLUMN "Review"."value_rating" IS 'Rating for value';
COMMENT ON COLUMN "Review"."accuracy_rating" IS 'Rating for listing accuracy';
COMMENT ON COLUMN "Review"."communication_rating" IS 'Rating for communication';
COMMENT ON COLUMN "Review"."cleanliness_rating" IS 'Rating for cleanliness';
COMMENT ON COLUMN "Review"."check_in_rating" IS 'Rating for check-in experience';

CREATE TABLE "SavedProperty" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "property_id" UUID NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  FOREIGN KEY ("property_id") REFERENCES "Property"("id") ON DELETE CASCADE,
  UNIQUE ("user_id", "property_id")
);

COMMENT ON TABLE "SavedProperty" IS 'Properties saved/favorited by users';
COMMENT ON COLUMN "SavedProperty"."user_id" IS 'Reference to the user (from auth microservice)';

CREATE TABLE "SearchHistory" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "search_query" TEXT NOT NULL,
  "filters" JSON,
  "created_at" TIMESTAMP DEFAULT now()
);

COMMENT ON TABLE "SearchHistory" IS 'User search history for properties';
COMMENT ON COLUMN "SearchHistory"."filters" IS 'Filters applied during search (location, price range, etc.)';

CREATE TABLE "PropertyNeighborhood" (
  "id" UUID PRIMARY KEY,
  "property_id" UUID NOT NULL,
  "school_district" VARCHAR,
  "school_rating" INT,
  "crime_rate" VARCHAR,
  "walk_score" INT,
  "transit_score" INT,
  "bike_score" INT,
  "nearby_locations" JSON,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  FOREIGN KEY ("property_id") REFERENCES "Property"("id") ON DELETE CASCADE
);

COMMENT ON TABLE "PropertyNeighborhood" IS 'Neighborhood information for properties';
COMMENT ON COLUMN "PropertyNeighborhood"."school_district" IS 'School district name';
COMMENT ON COLUMN "PropertyNeighborhood"."school_rating" IS 'School district rating (1-10)';
COMMENT ON COLUMN "PropertyNeighborhood"."crime_rate" IS 'Crime rate description';
COMMENT ON COLUMN "PropertyNeighborhood"."walk_score" IS 'Walkability score (0-100)';
COMMENT ON COLUMN "PropertyNeighborhood"."transit_score" IS 'Public transit score (0-100)';
COMMENT ON COLUMN "PropertyNeighborhood"."bike_score" IS 'Bikeability score (0-100)';
COMMENT ON COLUMN "PropertyNeighborhood"."nearby_locations" IS 'JSON array of nearby points of interest';

CREATE TABLE "PropertyView" (
  "id" UUID PRIMARY KEY,
  "property_id" UUID NOT NULL,
  "user_id" UUID NOT NULL,
  "viewed_at" TIMESTAMP DEFAULT now(),
  FOREIGN KEY ("property_id") REFERENCES "Property"("id") ON DELETE CASCADE
);

COMMENT ON TABLE "PropertyView" IS 'Tracks user views of properties';
COMMENT ON COLUMN "PropertyView"."id" IS 'Primary key for the PropertyView table';
COMMENT ON COLUMN "PropertyView"."property_id" IS 'Reference to the property viewed';
COMMENT ON COLUMN "PropertyView"."user_id" IS 'Reference to the user who viewed the property (from auth microservice)';
COMMENT ON COLUMN "PropertyView"."viewed_at" IS 'Timestamp of when the property was viewed';

-- Create indexes for performance
CREATE INDEX idx_property_category ON "Property"("category");
CREATE INDEX idx_property_type ON "Property"("type");
CREATE INDEX idx_property_status ON "Property"("status");
CREATE INDEX idx_property_city ON "Property"("city");
CREATE INDEX idx_property_state ON "Property"("state");
CREATE INDEX idx_property_country ON "Property"("country");
CREATE INDEX idx_property_price ON "Property"("price");
CREATE INDEX idx_property_detail_bedrooms ON "PropertyDetail"("bedrooms");
CREATE INDEX idx_property_detail_bathrooms ON "PropertyDetail"("bathrooms");
CREATE INDEX idx_property_detail_square_footage ON "PropertyDetail"("square_footage");
CREATE INDEX idx_property_detail_year_built ON "PropertyDetail"("year_built");
CREATE INDEX idx_amenity_category ON "Amenity"("category");
CREATE INDEX idx_review_property ON "Review"("property_id");
CREATE INDEX idx_review_overall_rating ON "Review"("overall_rating");
CREATE INDEX idx_availability_property_date ON "PropertyAvailability"("property_id", "date", "is_available");
CREATE INDEX idx_saved_property_user ON "SavedProperty"("user_id");
CREATE INDEX idx_search_history_user ON "SearchHistory"("user_id");
CREATE INDEX idx_property_neighborhood_property ON "PropertyNeighborhood"("property_id");
CREATE INDEX idx_property_view_user ON "PropertyView"("user_id");