package usecases

import (
	"context"
	"job_portal/property/internal/domain/entity"

	"job_portal/property/internal/domain/repository"

	"github.com/google/uuid"
)

type PropertyUsecase interface {
	CreateProperty(ctx context.Context, property *entity.Property) error
	GetProperties(ctx context.Context, limit, offset int32) ([]*entity.Property, error)
	GetPropertiesByOwner(ctx context.Context, ownerID uuid.NullUUID, limit, offset int32) ([]*entity.Property, error)
	GetPropertyByID(ctx context.Context, id uuid.UUID) (*entity.Property, error)
	UpdateProperty(ctx context.Context, property *entity.Property) error
	DeleteProperty(ctx context.Context, id uuid.UUID) error
}

type propertyUsecase struct {
	propertyRepo repository.PropertyRepository
}

func NewPropertyUsecase(propertyRepo repository.PropertyRepository) PropertyUsecase {
	return &propertyUsecase{
		propertyRepo: propertyRepo,
	}
}

// CreateProperty implements PropertyUsecase.
func (p *propertyUsecase) CreateProperty(ctx context.Context, property *entity.Property) error {
	if err := p.propertyRepo.CreateProperty(property); err != nil {
		return err
	}

	return nil
}

// DeleteProperty implements PropertyUsecase.
func (p *propertyUsecase) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	if err := p.propertyRepo.DeleteProperty(id); err != nil {
		return err
	}

	return nil
}

// GetProperties implements PropertyUsecase.
func (p *propertyUsecase) GetProperties(ctx context.Context, limit int32, offset int32) ([]*entity.Property, error) {
	property, err := p.propertyRepo.GetProperties(limit, offset)

	if err != nil {
		return nil, err
	}

	// get All Property
	properties := make([]*entity.Property, len(property))
	for i, property := range property {
		properties[i] = &entity.Property{
			ID:            property.ID,
			Title:         property.Title,
			Description:   property.Description,
			Price:         property.Price,
			Type:          string(property.Type),
			Address:       property.Address,
			ZipCode:       property.ZipCode,
			OwnerID:       property.OwnerID,
			Images:        property.Images,
			NoOfBedRooms:  property.NoOfBedRooms,
			NoOfBathRooms: property.NoOfBathRooms,
			NoOfToilets:   property.NoOfToilets,
			GeoLocation:   property.GeoLocation,
			Status:        string(property.Status),
			CreatedAt:     property.CreatedAt,
			UpdatedAt:     property.UpdatedAt,
		}
	}

	return properties, nil
}

// GetProperties implements PropertyUsecase.
func (p *propertyUsecase) GetPropertiesByOwner(ctx context.Context, ownerID uuid.NullUUID, limit int32, offset int32) ([]*entity.Property, error) {

	property, err := p.propertyRepo.GetPropertiesByOwner(ownerID, limit, offset)

	if err != nil {
		return nil, err
	}

	// get All Property

	properties := make([]*entity.Property, len(property))
	for i, property := range property {
		properties[i] = &entity.Property{
			ID:            property.ID,
			Title:         property.Title,
			Description:   property.Description,
			Price:         property.Price,
			Type:          string(property.Type),
			Address:       property.Address,
			ZipCode:       property.ZipCode,
			OwnerID:       property.OwnerID,
			Images:        property.Images,
			NoOfBedRooms:  property.NoOfBedRooms,
			NoOfBathRooms: property.NoOfBathRooms,
			NoOfToilets:   property.NoOfToilets,
			GeoLocation:   property.GeoLocation,
			Status:        string(property.Status),
			CreatedAt:     property.CreatedAt,
			UpdatedAt:     property.UpdatedAt,
		}
	}

	return properties, nil
}

// GetPropertyByID implements PropertyUsecase.
func (p *propertyUsecase) GetPropertyByID(ctx context.Context, id uuid.UUID) (*entity.Property, error) {

	property, err := p.propertyRepo.GetPropertyByID(id)

	if err != nil {
		return nil, err
	}

	return property, nil
}

// UpdateProperty implements PropertyUsecase.
func (p *propertyUsecase) UpdateProperty(ctx context.Context, property *entity.Property) error {
	if err := p.propertyRepo.UpdateProperty(property); err != nil {
		return err
	}

	return nil
}
