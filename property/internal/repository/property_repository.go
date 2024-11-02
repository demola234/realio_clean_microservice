package repository

import (
	"context"
	"database/sql"
	"fmt"
	db "job_portal/property/db/sqlc"
	"job_portal/property/internal/domain/entity"
	"strconv"

	"github.com/google/uuid"
)

// UserRepository implements the AuthRepository interface.
// This struct interacts with the database using SQLC-generated code.
type PropertyRepository struct {
	store db.Store
}

func NewPropertyRepository(store db.Store) *PropertyRepository {
	return &PropertyRepository{
		store: store,
	}
}

func (r *PropertyRepository) CreateProperty(property *entity.Property) error {

	NoOfBedRoom, err := strconv.Atoi(property.NoOfBedRooms)
	if err != nil {
		return err
	}

	NoOfBathRooms, err := strconv.Atoi(property.NoOfBathRooms)
	if err != nil {
		return err
	}

	NoOfToilets, err := strconv.Atoi(property.NoOfToilets)
	if err != nil {
		return err
	}

	arg := db.InsertPropertyParams{
		ID:            uuid.New(),
		Title:         property.Title,
		Description:   sql.NullString{String: property.Description, Valid: true},
		Price:         property.Price,
		Type:          db.NullPropertyType{PropertyType: db.PropertyType(property.Type), Valid: true},
		Address:       property.Address,
		ZipCode:       sql.NullString{String: property.ZipCode, Valid: true},
		OwnerID:       property.OwnerID,
		Images:        property.Images,
		NoOfBedRooms:  sql.NullInt32{Int32: int32(NoOfBedRoom), Valid: true},
		NoOfBathRooms: sql.NullInt32{Int32: int32(NoOfBathRooms), Valid: true},
		NoOfToilets:   sql.NullInt32{Int32: int32(NoOfToilets), Valid: true},
		GeoLocation:   property.GeoLocation,
		Status:        db.NullPropertyStatus{PropertyStatus: db.PropertyStatus(property.Status), Valid: true},
	}

	_, err = r.store.InsertProperty(context.Background(), arg)
	return err
}

func (r *PropertyRepository) UpdateProperty(property *entity.Property) error {
	NoOfBedRoom, err := strconv.Atoi(property.NoOfBedRooms)
	if err != nil {
		return err
	}

	NoOfBathRooms, err := strconv.Atoi(property.NoOfBathRooms)
	if err != nil {
		return err
	}

	NoOfToilets, err := strconv.Atoi(property.NoOfToilets)
	if err != nil {
		return err
	}

	arg := db.UpdatePropertyParams{
		ID:            property.ID,
		Title:         property.Title,
		Description:   sql.NullString{String: property.Description, Valid: true},
		Price:         property.Price,
		Type:          db.NullPropertyType{PropertyType: db.PropertyType(property.Type), Valid: true},
		Address:       property.Address,
		ZipCode:       sql.NullString{String: property.ZipCode, Valid: true},
		Images:        property.Images,
		NoOfBedRooms:  sql.NullInt32{Int32: int32(NoOfBedRoom), Valid: true},
		NoOfBathRooms: sql.NullInt32{Int32: int32(NoOfBathRooms), Valid: true},
		NoOfToilets:   sql.NullInt32{Int32: int32(NoOfToilets), Valid: true},
		GeoLocation:   property.GeoLocation,
		Status:        db.NullPropertyStatus{PropertyStatus: db.PropertyStatusAvailable, Valid: true},
	}

	return r.store.UpdateProperty(context.Background(), arg)
}

func (r *PropertyRepository) GetPropertyByID(id uuid.UUID) (*entity.Property, error) {

	property, err := r.store.GetPropertyByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &entity.Property{
		ID:            property.ID,
		Title:         property.Title,
		Description:   property.Description.String,
		Price:         property.Price,
		Type:          string(property.Type.PropertyType),
		Address:       property.Address,
		ZipCode:       property.ZipCode.String,
		OwnerID:       property.OwnerID,
		Images:        property.Images,
		NoOfBedRooms:  strconv.Itoa(int(property.NoOfBedRooms.Int32)),
		NoOfBathRooms: strconv.Itoa(int(property.NoOfBathRooms.Int32)),
		NoOfToilets:   strconv.Itoa(int(property.NoOfToilets.Int32)),
		GeoLocation:   property.GeoLocation,
		Status:        string(property.Status.PropertyStatus),
		CreatedAt:     property.CreatedAt.Time,
		UpdatedAt:     property.UpdatedAt.Time,
	}, nil

}

func (r *PropertyRepository) GetProperties(limit, offset int32) ([]*entity.Property, error) {
	properties, err := r.store.ListProperties(context.Background(), db.ListPropertiesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to Get Properties: %w", err)
	}

	result := make([]*entity.Property, len(properties))
	for i, property := range properties {
		result[i] = &entity.Property{
			ID:            property.ID,
			Title:         property.Title,
			Description:   property.Description.String,
			Price:         property.Price,
			Type:          string(property.Type.PropertyType),
			Address:       property.Address,
			ZipCode:       property.ZipCode.String,
			OwnerID:       property.OwnerID,
			Images:        property.Images,
			NoOfBedRooms:  strconv.Itoa(int(property.NoOfBedRooms.Int32)),
			NoOfBathRooms: strconv.Itoa(int(property.NoOfBathRooms.Int32)),
			NoOfToilets:   strconv.Itoa(int(property.NoOfToilets.Int32)),
			GeoLocation:   property.GeoLocation,
			Status:        string(property.Status.PropertyStatus),
			CreatedAt:     property.CreatedAt.Time,
			UpdatedAt:     property.UpdatedAt.Time,
		}
	}

	return result, nil
}

func (r *PropertyRepository) GetPropertiesByOwner(ownerID uuid.NullUUID, limit, offset int32) ([]*entity.Property, error) {
	properties, err := r.store.GetPropertiesByOwnerID(context.Background(), db.GetPropertiesByOwnerIDParams{
		OwnerID: ownerID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Property, len(properties))
	for i, property := range properties {
		result[i] = &entity.Property{
			ID:            property.ID,
			Title:         property.Title,
			Description:   property.Description.String,
			Price:         property.Price,
			Type:          string(property.Type.PropertyType),
			Address:       property.Address,
			ZipCode:       property.ZipCode.String,
			OwnerID:       property.OwnerID,
			Images:        property.Images,
			NoOfBedRooms:  strconv.Itoa(int(property.NoOfBedRooms.Int32)),
			NoOfBathRooms: strconv.Itoa(int(property.NoOfBathRooms.Int32)),
			NoOfToilets:   strconv.Itoa(int(property.NoOfToilets.Int32)),
			GeoLocation:   property.GeoLocation,
			Status:        string(property.Status.PropertyStatus),
			CreatedAt:     property.CreatedAt.Time,
			UpdatedAt:     property.UpdatedAt.Time,
		}
	}

	return result, nil
}

func (r *PropertyRepository) DeleteProperty(id uuid.UUID) error {
	return r.store.DeleteProperty(context.Background(), id)
}
