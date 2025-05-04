package repository

import (
	"github.com/demola234/property/internal/domain/entity"

	"github.com/google/uuid"
)

type PropertyRepository interface {
	CreateProperty(property *entity.Property) error
	GetProperties(limit, offset int32) ([]*entity.Property, error)
	GetPropertyByID(id uuid.UUID) (*entity.Property, error)
	UpdateProperty(property *entity.Property) error
	DeleteProperty(id uuid.UUID) error
	GetPropertiesByOwner(ownerID uuid.NullUUID, limit, offset int32) ([]*entity.Property, error)
}
