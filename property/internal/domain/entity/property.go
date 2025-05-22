package entity

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
)

// Enums
type PropertyCategory string

const (
	PropertyCategoryRent     PropertyCategory = "Rent"
	PropertyCategorySale     PropertyCategory = "Sale"
	PropertyCategoryBuy      PropertyCategory = "Buy"
	PropertyCategoryLease    PropertyCategory = "Lease"
	PropertyCategoryLand     PropertyCategory = "Land"
	PropertyCategoryMortgage PropertyCategory = "Mortgage"
)

type PropertyType string

const (
	PropertyTypeHouse        PropertyType = "House"
	PropertyTypeApartment    PropertyType = "Apartment"
	PropertyTypeLand         PropertyType = "Land"
	PropertyTypeManufactured PropertyType = "Manufactured"
	PropertyTypeTownhome     PropertyType = "Townhome"
	PropertyTypeMultiFamily  PropertyType = "Multi-family"
	PropertyTypeCondo        PropertyType = "Condo"
	PropertyTypeCoop         PropertyType = "Co-op"
	PropertyTypeLot          PropertyType = "Lot"
)

type PropertyStatus string

const (
	PropertyStatusAvailable PropertyStatus = "Available"
	PropertyStatusSold      PropertyStatus = "Sold"
	PropertyStatusRented    PropertyStatus = "Rented"
	PropertyStatusPending   PropertyStatus = "Pending"
)

type AmenityCategory string

const (
	AmenityCategoryFeature   AmenityCategory = "Feature"
	AmenityCategoryAmenity   AmenityCategory = "Amenity"
	AmenityCategoryAppliance AmenityCategory = "Appliance"
	AmenityCategoryUtility   AmenityCategory = "Utility"
	AmenityCategoryParking   AmenityCategory = "Parking"
)

type ReviewType string

const (
	ReviewTypeProperty ReviewType = "Property"
	ReviewTypeHost     ReviewType = "Host"
	ReviewTypeGuest    ReviewType = "Guest"
)

type Property struct {
	ID          uuid.UUID        `json:"id"`
	Title       string           `json:"title"`
	Description *string          `json:"description"`
	Price       decimal.Decimal  `json:"price"`
	Category    PropertyCategory `json:"category"`
	Type        PropertyType     `json:"type"`
	Address     string           `json:"address"`
	City        string           `json:"city"`
	State       string           `json:"state"`
	Country     string           `json:"country"`
	ZipCode     *string          `json:"zip_code"`
	OwnerID     *uuid.UUID       `json:"owner_id"`
	Status      PropertyStatus   `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type PropertyDetail struct {
	ID               uuid.UUID              `json:"id"`
	PropertyID       uuid.UUID              `json:"property_id"`
	Bedrooms         *int32                 `json:"bedrooms"`
	Bathrooms        *int32                 `json:"bathrooms"`
	Toilets          *int32                 `json:"toilets"`
	SquareFootage    *decimal.Decimal       `json:"square_footage"`
	LotSize          *decimal.Decimal       `json:"lot_size"`
	YearBuilt        *int32                 `json:"year_built"`
	Stories          *int32                 `json:"stories"`
	GarageCount      *int32                 `json:"garage_count"`
	HasBasement      bool                   `json:"has_basement"`
	HasAttic         bool                   `json:"has_attic"`
	HeatingSystem    *string                `json:"heating_system"`
	CoolingSystem    *string                `json:"cooling_system"`
	WaterSource      *string                `json:"water_source"`
	SewerType        *string                `json:"sewer_type"`
	RoofType         *string                `json:"roof_type"`
	ExteriorMaterial *string                `json:"exterior_material"`
	FoundationType   *string                `json:"foundation_type"`
	PoolType         *string                `json:"pool_type"`
	GeoLocation      map[string]interface{} `json:"geo_location"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type Amenity struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Icon        *string         `json:"icon"`
	Category    AmenityCategory `json:"category"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type PropertyAmenity struct {
	PropertyID uuid.UUID `json:"property_id"`
	AmenityID  uuid.UUID `json:"amenity_id"`
	HasAmenity bool      `json:"has_amenity"`
	Notes      *string   `json:"notes"`
}

type PropertyImage struct {
	ID           uuid.UUID `json:"id"`
	PropertyID   uuid.UUID `json:"property_id"`
	URL          string    `json:"url"`
	Caption      *string   `json:"caption"`
	IsPrimary    bool      `json:"is_primary"`
	DisplayOrder int32     `json:"display_order"`
	RoomType     *string   `json:"room_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Review struct {
	ID                  uuid.UUID        `json:"id"`
	BookingID           *uuid.UUID       `json:"booking_id"`
	PropertyID          uuid.UUID        `json:"property_id"`
	ReviewerID          uuid.UUID        `json:"reviewer_id"`
	ReviewedID          *uuid.UUID       `json:"reviewed_id"`
	OverallRating       decimal.Decimal  `json:"overall_rating"`
	LocationRating      *decimal.Decimal `json:"location_rating"`
	ValueRating         *decimal.Decimal `json:"value_rating"`
	AccuracyRating      *decimal.Decimal `json:"accuracy_rating"`
	CommunicationRating *decimal.Decimal `json:"communication_rating"`
	CleanlinessRating   *decimal.Decimal `json:"cleanliness_rating"`
	CheckInRating       *decimal.Decimal `json:"check_in_rating"`
	Comment             *string          `json:"comment"`
	Type                ReviewType       `json:"type"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

// Search and Filter Types
type PropertySearchFilter struct {
	Category        *PropertyCategory `json:"category"`
	Type            *PropertyType     `json:"type"`
	Status          *PropertyStatus   `json:"status"`
	City            *string           `json:"city"`
	State           *string           `json:"state"`
	Country         *string           `json:"country"`
	MinPrice        *decimal.Decimal  `json:"min_price"`
	MaxPrice        *decimal.Decimal  `json:"max_price"`
	MinBedrooms     *int32            `json:"min_bedrooms"`
	MinBathrooms    *int32            `json:"min_bathrooms"`
	MinSquareFeet   *int32            `json:"min_square_feet"`
	MinYearBuilt    *int32            `json:"min_year_built"`
	MinGarageCount  *int32            `json:"min_garage_count"`
	HasBasement     *bool             `json:"has_basement"`
	HasAttic        *bool             `json:"has_attic"`
	MinWalkScore    *int32            `json:"min_walk_score"`
	MinSchoolRating *int32            `json:"min_school_rating"`
	SortBy          string            `json:"sort_by"`
	Limit           int32             `json:"limit"`
	Offset          int32             `json:"offset"`
}

type PropertyWithDetails struct {
	Property
	PropertyDetail *PropertyDetail `json:"property_detail"`
	AvgRating   decimal.Decimal `json:"avg_rating"`
	ReviewCount int64           `json:"review_count"`
}
