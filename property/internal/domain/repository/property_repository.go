package repository

import (
	"context"

	"github.com/demola234/property/internal/domain/entity"

	"github.com/google/uuid"
)

type PropertyRepository interface {
	// Property CRUD
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Property, error)
	GetWithAllDetails(ctx context.Context, id uuid.UUID) (*entity.PropertyWithDetails, error)
	Create(ctx context.Context, property *entity.Property) error
	Update(ctx context.Context, property *entity.Property) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.PropertyStatus) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Property Listing and Search
	List(ctx context.Context, filter entity.PropertySearchFilter) ([]*entity.Property, error)
	SearchWithDetails(ctx context.Context, filter entity.PropertySearchFilter) ([]*entity.PropertyWithDetails, error)
	Count(ctx context.Context, filter entity.PropertySearchFilter) (int64, error)
	GetByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int32) ([]*entity.Property, error)

	// Property with amenities
	GetByAmenity(ctx context.Context, amenityID uuid.UUID, status *entity.PropertyStatus, limit, offset int32) ([]*entity.Property, error)
	GetByMultipleAmenities(ctx context.Context, amenityIDs []uuid.UUID, status *entity.PropertyStatus, limit, offset int32) ([]*entity.Property, error)
}
