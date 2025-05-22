package repository

import (
	"context"

	"github.com/demola234/property/internal/domain/entity"

	"github.com/google/uuid"
)

type PropertyAmenityRepository interface {
	GetPropertyAmenities(ctx context.Context, propertyID uuid.UUID) ([]*entity.PropertyAmenity, error)
	AddPropertyAmenity(ctx context.Context, propertyAmenity *entity.PropertyAmenity) error
	RemovePropertyAmenity(ctx context.Context, propertyID, amenityID uuid.UUID) error
	RemoveAllPropertyAmenities(ctx context.Context, propertyID uuid.UUID) error
}