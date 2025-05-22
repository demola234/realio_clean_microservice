package repository

import (
	"context"

	"github.com/demola234/property/internal/domain/entity"

	"github.com/google/uuid"
)

type PropertyViewRepository interface {
	RecordView(ctx context.Context, propertyID, userID uuid.UUID) error
	GetViewStats(ctx context.Context, propertyID uuid.UUID) (*entity.ViewStats, error)
	GetRecentlyViewed(ctx context.Context, userID uuid.UUID, limit int32) ([]*entity.Property, error)
}
