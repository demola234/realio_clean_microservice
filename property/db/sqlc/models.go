// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type PropertyStatus string

const (
	PropertyStatusAvailable PropertyStatus = "Available"
	PropertyStatusSold      PropertyStatus = "Sold"
	PropertyStatusRented    PropertyStatus = "Rented"
)

func (e *PropertyStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PropertyStatus(s)
	case string:
		*e = PropertyStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PropertyStatus: %T", src)
	}
	return nil
}

type NullPropertyStatus struct {
	PropertyStatus PropertyStatus `json:"property_status"`
	Valid          bool           `json:"valid"` // Valid is true if PropertyStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPropertyStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PropertyStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PropertyStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPropertyStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PropertyStatus), nil
}

type PropertyType string

const (
	PropertyTypeHouse     PropertyType = "House"
	PropertyTypeApartment PropertyType = "Apartment"
	PropertyTypeLand      PropertyType = "Land"
)

func (e *PropertyType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PropertyType(s)
	case string:
		*e = PropertyType(s)
	default:
		return fmt.Errorf("unsupported scan type for PropertyType: %T", src)
	}
	return nil
}

type NullPropertyType struct {
	PropertyType PropertyType `json:"property_type"`
	Valid        bool         `json:"valid"` // Valid is true if PropertyType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPropertyType) Scan(value interface{}) error {
	if value == nil {
		ns.PropertyType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PropertyType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPropertyType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PropertyType), nil
}

type Property struct {
	// Primary key
	ID uuid.UUID `json:"id"`
	// Property title
	Title string `json:"title"`
	// Detailed description
	Description sql.NullString `json:"description"`
	// Price of the property
	Price string `json:"price"`
	// House, Apartment, Land, etc.
	Type NullPropertyType `json:"type"`
	// Address details
	Address string `json:"address"`
	// Zip code (Optional)
	ZipCode sql.NullString `json:"zip_code"`
	// Reference to the user (seller), external microservice
	OwnerID uuid.NullUUID `json:"owner_id"`
	// List of image URLs (Optional)
	Images []string `json:"images"`
	// Number of Bed Rooms
	NoOfBedRooms sql.NullInt32 `json:"no_of_bed_rooms"`
	// Number of Bath Rooms
	NoOfBathRooms sql.NullInt32 `json:"no_of_bath_rooms"`
	// Number of Toilets
	NoOfToilets sql.NullInt32 `json:"no_of_toilets"`
	// Latitude & longitude (Optional)
	GeoLocation pqtype.NullRawMessage `json:"geo_location"`
	// Available, Sold, Rented, etc.
	Status NullPropertyStatus `json:"status"`
	// Timestamp of listing
	CreatedAt sql.NullTime `json:"created_at"`
	// Timestamp of last update
	UpdatedAt sql.NullTime `json:"updated_at"`
}
