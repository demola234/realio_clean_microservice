-- Step 1: Create ENUM types for `type` and `status`
CREATE TYPE property_type AS ENUM ('House', 'Apartment', 'Land');
CREATE TYPE property_status AS ENUM ('Available', 'Sold', 'Rented');

-- Step 2: Create the `Property` table with comments
CREATE TABLE "Property" (
  "id" UUID PRIMARY KEY,
  "title" VARCHAR NOT NULL,
  "description" TEXT,
  "price" NUMERIC NOT NULL,
  "type" property_type, -- Enum for property type
  "address" VARCHAR NOT NULL,
  "zip_code" VARCHAR,
  "owner_id" UUID, -- Reference to external user microservice
  "images" TEXT[], -- Array of image URLs (Optional)
  "no_of_bed_rooms" INT,
  "no_of_bath_rooms" INT,
  "no_of_toilets" INT,
  "geo_location" JSON, -- Latitude & longitude (Optional)
  "status" property_status, -- Enum for property status
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

-- Step 3: Add comments for each column
COMMENT ON COLUMN "Property"."id" IS 'Primary key';
COMMENT ON COLUMN "Property"."title" IS 'Property title';
COMMENT ON COLUMN "Property"."description" IS 'Detailed description';
COMMENT ON COLUMN "Property"."price" IS 'Price of the property';
COMMENT ON COLUMN "Property"."type" IS 'House, Apartment, Land, etc.';
COMMENT ON COLUMN "Property"."address" IS 'Address details';
COMMENT ON COLUMN "Property"."zip_code" IS 'Zip code (Optional)';
COMMENT ON COLUMN "Property"."owner_id" IS 'Reference to the user (seller), external microservice';
COMMENT ON COLUMN "Property"."images" IS 'List of image URLs (Optional)';
COMMENT ON COLUMN "Property"."no_of_bed_rooms" IS 'Number of Bed Rooms';
COMMENT ON COLUMN "Property"."no_of_bath_rooms" IS 'Number of Bath Rooms';
COMMENT ON COLUMN "Property"."no_of_toilets" IS 'Number of Toilets';
COMMENT ON COLUMN "Property"."geo_location" IS 'Latitude & longitude (Optional)';
COMMENT ON COLUMN "Property"."status" IS 'Available, Sold, Rented, etc.';
COMMENT ON COLUMN "Property"."created_at" IS 'Timestamp of listing';
COMMENT ON COLUMN "Property"."updated_at" IS 'Timestamp of last update';