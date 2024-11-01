CREATE TABLE "Property" (
  "id" UUID PRIMARY KEY,
  "title" String NOT NULL,
  "description" Text,
  "price" Decimal NOT NULL,
  "type" Enum,
  "address" String NOT NULL,
  "zip_code" String,
  "owner_id" UUID,
  "images" Array,
  "no_of_bed_rooms" Int,
  "no_of_bath_rooms" Int,
  "no_of_toilets" Int,
  "geo_location" JSON,
  "status" Enum,
  "created_at" Timestamp DEFAULT (now()),
  "updated_at" Timestamp DEFAULT (now())
);

COMMENT ON COLUMN "Property"."id" IS 'Primary key';

COMMENT ON COLUMN "Property"."title" IS 'Property title';

COMMENT ON COLUMN "Property"."description" IS 'Detailed description';

COMMENT ON COLUMN "Property"."price" IS 'Price of the property';

COMMENT ON COLUMN "Property"."type" IS 'House, Apartment, Land, etc. (Optional)';

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
