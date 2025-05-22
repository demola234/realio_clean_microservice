package repository

import (
	"context"

	"github.com/demola234/property/internal/domain/entity"

	"github.com/google/uuid"
)

type ReviewRepository interface {
	GetPropertyReviews(ctx context.Context, propertyID uuid.UUID, limit, offset int32) ([]*entity.Review, error)
	GetPropertyReviewStats(ctx context.Context, propertyID uuid.UUID) (*entity.ReviewStats, error)
	Create(ctx context.Context, review *entity.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
}
