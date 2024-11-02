package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Property struct {
	ID            uuid.UUID             `json:"id"`
	Title         string                `json:"title"`
	Description   string                `json:"description"`
	Price         string                `json:"price"`
	Type          string                `json:"type"`
	Address       string                `json:"address"`
	ZipCode       string                `json:"zip_code"`
	OwnerID       uuid.NullUUID         `json:"owner_id"`
	Images        []string              `json:"images"`
	NoOfBedRooms  string                `json:"no_of_bed_rooms"`
	NoOfBathRooms string                `json:"no_of_bath_rooms"`
	NoOfToilets   string                `json:"no_of_toilets"`
	GeoLocation   pqtype.NullRawMessage `json:"geo_location"`
	Status        string                `json:"status"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}
