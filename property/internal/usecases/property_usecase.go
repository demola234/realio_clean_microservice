package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"job_portal/property/internal/domain/entity"
	"job_portal/property/internal/domain/repository"
	"log"

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
	messageQueue entity.MessageQueue
}

func NewPropertyUsecase(propertyRepo repository.PropertyRepository, mq entity.MessageQueue) PropertyUsecase {
	return &propertyUsecase{
		propertyRepo: propertyRepo,
		messageQueue: mq,
	}
}

// CreateProperty creates a new property in the repository.
func (p *propertyUsecase) CreateProperty(ctx context.Context, property *entity.Property) error {
	if err := p.propertyRepo.CreateProperty(property); err != nil {
		return fmt.Errorf("failed to create property: %w", err)
	}

	// Prepare the event data to publish to Kafka
	eventData, err := json.Marshal(property)
	if err != nil {
		return fmt.Errorf("failed to marshal property data for event: %w", err)
	}

	// Publish the event to Kafka
	if err := p.messageQueue.PublishMessage(ctx, []byte("property_created"), eventData); err != nil {
		return fmt.Errorf("failed to publish property created event: %w", err)
	}

	return nil
}

// DeleteProperty deletes a property from the repository.
func (p *propertyUsecase) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	if err := p.propertyRepo.DeleteProperty(id); err != nil {
		return fmt.Errorf("failed to delete property with ID %s: %w", id, err)
	}

	// Step 2: Publish a delete event to Kafka
	eventData, err := json.Marshal(map[string]interface{}{
		"event": "property_deleted",
		"id":    id.String(),
	})
	if err != nil {
		log.Printf("Failed to marshal delete event data: %v", err)
		return err
	}

	if err := p.messageQueue.PublishMessage(ctx, []byte("property_deleted"), eventData); err != nil {
		log.Printf("Failed to publish property deleted event: %v", err)
		return err
	}

	return nil
}

// GetProperties retrieves a list of properties with pagination.
func (p *propertyUsecase) GetProperties(ctx context.Context, limit, offset int32) ([]*entity.Property, error) {
	
	return p.propertyRepo.GetProperties(limit, offset)
}

// GetPropertiesByOwner retrieves a list of properties owned by a specific owner.
func (p *propertyUsecase) GetPropertiesByOwner(ctx context.Context, ownerID uuid.NullUUID, limit, offset int32) ([]*entity.Property, error) {
	
	return p.propertyRepo.GetPropertiesByOwner(ownerID, limit, offset)
}

// GetPropertyByID retrieves a single property by its ID.
func (p *propertyUsecase) GetPropertyByID(ctx context.Context, id uuid.UUID) (*entity.Property, error) {

	return p.propertyRepo.GetPropertyByID(id)
}

// UpdateProperty updates an existing property in the repository.
func (p *propertyUsecase) UpdateProperty(ctx context.Context, property *entity.Property) error {

	return p.propertyRepo.UpdateProperty(property)
}
