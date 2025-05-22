package repository

import (
	"context"

	"github.com/demola234/property/internal/domain/entity"

	"github.com/google/uuid"
)

type PropertyDetailRepository interface {
	GetByPropertyID(ctx context.Context, propertyID uuid.UUID) (*entity.PropertyDetail, error)
	Create(ctx context.Context, detail *entity.PropertyDetail) error
	Update(ctx context.Context, detail *entity.PropertyDetail) error
}
