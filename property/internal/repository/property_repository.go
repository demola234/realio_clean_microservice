package repository

import (
	"context"
	"database/sql"

	db "github.com/demola234/property/db/sqlc"
	"github.com/demola234/property/internal/domain/entity"

	"github.com/google/uuid"
)

type PropertyRepository struct {
	store db.Store
}

func NewPropertyRepository(store db.Store) *PropertyRepository {
	return &PropertyRepository{
		store: store,
	}
}

func (r *PropertyRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Property, error) {
	property, err := r.store.GetPropertyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.mapSQLCPropertyToEntity(property), nil
}

func (r *PropertyRepository) GetWithAllDetails(ctx context.Context, id uuid.UUID) (*entity.PropertyWithDetails, error) {
	result, err := r.store.GetPropertyWithAllDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.mapSQLCPropertyWithDetailsToEntity(result), nil
}

func (r *PropertyRepository) Create(ctx context.Context, property *entity.Property) error {
	params := db.CreatePropertyParams{
		Title:       property.Title,
		Description: sql.NullString{String: *property.Description, Valid: property.Description != nil},
		Price:       property.Price.String(),
		Category:    db.PropertyCategory(property.Category),
		Type:        db.PropertyType(property.Type),
		Address:     property.Address,
		City:        property.City,
		State:       property.State,
		Country:     property.Country,
		ZipCode:     sql.NullString{String: *property.ZipCode, Valid: property.ZipCode != nil},
		OwnerID:     toNullUUID(property.OwnerID),
		Status:      db.PropertyStatus(property.Status),
	}

	created, err := r.store.CreateProperty(ctx, params)
	if err != nil {
		return err
	}

	property.ID = created.ID
	property.CreatedAt = toTime(created.CreatedAt)
	property.UpdatedAt = toTime(created.UpdatedAt)
	return nil
}

func (r *PropertyRepository) Update(ctx context.Context, property *entity.Property) error {
	params := db.UpdatePropertyParams{
		ID:          property.ID,
		Title:       property.Title,
		Description: sql.NullString{String: *property.Description, Valid: property.Description != nil},
		Category:    db.PropertyCategory(property.Category),
		Type:        db.PropertyType(property.Type),
		Address:     property.Address,
		City:        property.City,
		State:       property.State,
		Country:     property.Country,
		ZipCode:     sql.NullString{String: *property.ZipCode, Valid: property.ZipCode != nil},
		Status:      db.PropertyStatus(property.Status),
	}

	_, err := r.store.UpdateProperty(ctx, params)
	if err != nil {
		return err
	}

	property.UpdatedAt = sql.NullTime{Time: property.UpdatedAt, Valid: true}.Time
	return nil
}

func (r *PropertyRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.PropertyStatus) error {
	_, err := r.store.UpdatePropertyStatus(ctx, db.UpdatePropertyStatusParams{
		ID:     id,
		Status: db.PropertyStatus(status),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *PropertyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteProperty(ctx, id)
	return err
}

func (r *PropertyRepository) List(ctx context.Context, filter entity.PropertySearchFilter) ([]*entity.Property, error) {
	params := db.ListPropertiesParams{
		Category:     (*db.PropertyCategory)(filter.Category),
		Type:         (*db.PropertyType)(filter.Type),
		Status:       (*db.PropertyStatus)(filter.Status),
		City:         sql.NullString{String: *filter.City, Valid: filter.City != nil},
		State:        sql.NullString{String: *filter.State, Valid: filter.State != nil},
		Country:      sql.NullString{String: *filter.Country, Valid: filter.Country != nil},
		MinPrice:     sql.NullString{String: filter.MinPrice, Valid: filter.MinPrice != nil},
		MaxPrice:     sql.NullString{String: filter.MaxPrice, Valid: filter.MaxPrice != nil},
		MinBedrooms:  sql.NullInt32{Int32: *filter.MinBedrooms, Valid: filter.MinBedrooms != nil},
		MinBathrooms: sql.NullInt32{Int32: *filter.MinBathrooms, Valid: filter.MinBathrooms != nil},
		Limit:        filter.Limit,
		Offset:       filter.Offset,
	}

	properties, err := r.store.ListProperties(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Property, len(properties))
	for i, p := range properties {
		result[i] = r.mapSQLCListPropertyToEntity(p)
	}

	return result, nil
}

// GetWithAllDetails(ctx context.Context, id uuid.UUID) (*entity.PropertyWithDetails, error)
// Create(ctx context.Context, property *entity.Property) error
// Update(ctx context.Context, property *entity.Property) error
// UpdateStatus(ctx context.Context, id uuid.UUID, status entity.PropertyStatus) error
// Delete(ctx context.Context, id uuid.UUID) error

// // Property Listing and Search
// List(ctx context.Context, filter entity.PropertySearchFilter) ([]*entity.Property, error)
// SearchWithDetails(ctx context.Context, filter entity.PropertySearchFilter) ([]*entity.PropertyWithDetails, error)
// Count(ctx context.Context, filter entity.PropertySearchFilter) (int64, error)
// GetByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int32) ([]*entity.Property, error)

// // Property with amenities
// GetByAmenity(ctx context.Context, amenityID uuid.UUID, status *entity.PropertyStatus, limit, offset int32) ([]*entity.Property, error)
// GetByMultipleAmenities(ctx context.Context, amenityIDs []uuid.UUID, status *entity.PropertyStatus, limit, offset int32) ([]*entity.Property, error)

func (r *PropertyRepository) mapSQLCPropertyToEntity(p db.Property) *entity.Property {
	return &entity.Property{
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

func (r *PropertyRepository) mapSQLCListPropertyToEntity(p db.ListPropertiesRow) *entity.Property {
	return &entity.Property{
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

func (r *PropertyRepository) mapSQLCPropertyWithDetailsToEntity(p db.GetPropertyWithAllDetailsRow) *entity.PropertyWithDetails {
	return &entity.PropertyWithDetails{
		Property: entity.Property{
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
		},
		PropertyDetail: &entity.PropertyDetail{
			ID:               p.ID_2.UUID,
			PropertyID:       p.PropertyID.UUID,
			Bedrooms:         toNullInt32Ptr(p.Bedrooms),
			Bathrooms:        toNullInt32Ptr(p.Bathrooms),
			Toilets:          toNullInt32Ptr(p.Toilets),
			SquareFootage:    toProtoDecimal(p.SquareFootage),
			LotSize:          toProtoDecimal(p.LotSize),
			YearBuilt:        toNullInt32Ptr(p.YearBuilt),
			Stories:          toNullInt32Ptr(p.Stories),
			GarageCount:      toNullInt32Ptr(p.GarageCount),
			HasBasement:      p.HasBasement.Bool,
			HasAttic:         p.HasAttic.Bool,
			HeatingSystem:    toNullStringPtr(p.HeatingSystem),
			CoolingSystem:    toNullStringPtr(p.CoolingSystem),
			WaterSource:      toNullStringPtr(p.WaterSource),
			SewerType:        toNullStringPtr(p.SewerType),
			RoofType:         toNullStringPtr(p.RoofType),
			ExteriorMaterial: toNullStringPtr(p.ExteriorMaterial),
			FoundationType:   toNullStringPtr(p.FoundationType),
			PoolType:         toNullStringPtr(p.PoolType),
			CreatedAt:        p.CreatedAt_2.Time,
			UpdatedAt:        p.UpdatedAt_2.Time,
		},
		AvgRating:   *float64ToProtoDecimal(p.AvgRating),
		ReviewCount: p.ReviewCount,
	}
}

func (r *PropertyRepository) mapSQLCSearchResultToEntity(p db.SearchPropertiesWithDetailsRow) *entity.PropertyWithDetails {
	return &entity.PropertyWithDetails{
		Property: entity.Property{
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
		},
		PropertyDetail: &entity.PropertyDetail{
			Bedrooms:      toNullInt32Ptr(p.Bedrooms),
			Bathrooms:     toNullInt32Ptr(p.Bathrooms),
			SquareFootage: toProtoDecimal(p.SquareFootage),
			LotSize:       toProtoDecimal(p.LotSize),
			YearBuilt:     toNullInt32Ptr(p.YearBuilt),
			Stories:       toNullInt32Ptr(p.Stories),
			GarageCount:   toNullInt32Ptr(p.GarageCount),
			HasBasement:   p.HasBasement.Bool,
			HasAttic:      p.HasAttic.Bool,
		},
		AvgRating:   *float64ToProtoDecimal(p.AvgRating),
		ReviewCount: p.ReviewCount,
	}
}
