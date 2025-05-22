package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/demola234/property/db/sqlc"
	"github.com/demola234/property/internal/domain/entity"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type PropertyViewRepository struct {
	store db.Store
}

func NewPropertyViewRepository(store db.Store) *PropertyViewRepository {
	return &PropertyViewRepository{
		store: store,
	}
}

func (r *PropertyViewRepository) RecordView(ctx context.Context, propertyID, userID uuid.UUID) error {

	params := db.RecordPropertyViewParams{
		PropertyID: propertyID,
		UserID:     userID,
	}

	_, err := r.store.RecordPropertyView(context.Background(), params)

	return err
}

func (r *PropertyViewRepository) GetViewStats(ctx context.Context, propertyID uuid.UUID) (*entity.ViewStats, error) {
	stats, err := r.store.GetPropertyViewStats(ctx, propertyID)
	if err != nil {
		return nil, err
	}

	return &entity.ViewStats{
		TotalViews:  stats.TotalViews,
		UniqueViews: stats.UniqueViews,
	}, nil
}

func (r *PropertyViewRepository) GetRecentlyViewed(ctx context.Context, userID uuid.UUID, limit int32) ([]*entity.Property, error) {
	params := db.GetRecentlyViewedPropertiesParams{
		UserID: userID,
		Limit:  limit,
	}

	properties, err := r.store.GetRecentlyViewedProperties(ctx, params)

	if err != nil {
		return nil, err
	}

	result := make([]*entity.Property, len(properties))
	for i, p := range properties {
		result[i] = &entity.Property{
			ID:          p.ID,
			Title:       p.Title,
			Description: toNullStringPtr(p.Description),
			Category:    entity.PropertyCategory(p.Category),
			Type:        entity.PropertyType(p.Type),
			Address:     p.Address,
			City:        p.City,
			State:       p.State,
			Country:     p.Country,
			ZipCode:     toNullStringPtr(p.ZipCode),
			OwnerID:     toNullUUIDPtr(p.OwnerID),
			Status:      entity.PropertyStatus(p.Status),
			CreatedAt:   toTime(p.CreatedAt),
			UpdatedAt:   toTime(p.UpdatedAt),
		}
	}

	return result, nil

}

func toNullStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func toNullInt32Ptr(ni sql.NullInt32) *int32 {
	if ni.Valid {
		return &ni.Int32
	}
	return nil
}

func toNullInt64Ptr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

func toProtoDecimal(ns sql.NullString) *decimal.Decimal {
	if !ns.Valid {
		return nil
	}
	return &decimal.Decimal{
		Value: ns.String,
	}
}

func toNullUUIDPtr(nu uuid.NullUUID) *uuid.UUID {
	if nu.Valid {
		return &nu.UUID
	}
	return nil
}

func toTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

func toTimePtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func float64ToProtoDecimal(f float64) *decimal.Decimal {
	return &decimal.Decimal{
		Value: fmt.Sprintf("%.6f", f),
	}
}

func toNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u != nil {
		return uuid.NullUUID{UUID: *u, Valid: true}
	}
	return uuid.NullUUID{Valid: false}
}
