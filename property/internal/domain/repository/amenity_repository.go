package repository

import (
	"context"

	"github.com/demola234/property/internal/domain/entity"

	"github.com/google/uuid"
)

type AmenityRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Amenity, error)
	List(ctx context.Context, category *entity.AmenityCategory) ([]*entity.Amenity, error)
	Create(ctx context.Context, amenity *entity.Amenity) error
	Update(ctx context.Context, amenity *entity.Amenity) error
	Delete(ctx context.Context, id uuid.UUID) error
}
