package repository

import (
	"context"

	"github.com/demola234/property/internal/domain/entity"

	"github.com/google/uuid"
)

type PropertyImageRepository interface {
	ListByProperty(ctx context.Context, propertyID uuid.UUID) ([]*entity.PropertyImage, error)
	Add(ctx context.Context, image *entity.PropertyImage) error
	Update(ctx context.Context, image *entity.PropertyImage) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetPrimary(ctx context.Context, propertyID, imageID uuid.UUID) error
}
