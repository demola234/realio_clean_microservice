package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/demola234/property/pkg/utils"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/require"
)

func createRandomProperty(t *testing.T) Property {
	street := fmt.Sprintf("%d %s", utils.RandomInt(1, 100), utils.RandomString(6))
	arg := InsertPropertyParams{
		ID:            uuid.New(),
		Title:         utils.RandomString(6),
		Description:   sql.NullString{String: utils.RandomString(20), Valid: true},
		Price:         strconv.Itoa(utils.RandomInt(100000, 1000000)),
		Type:          NullPropertyType{PropertyType: "House", Valid: true},
		Address:       street,
		ZipCode:       sql.NullString{String: "12345", Valid: true},
		OwnerID:       uuid.NullUUID{UUID: uuid.New(), Valid: true},
		Images:        []string{"image1.jpg", "image2.jpg"},
		NoOfBedRooms:  sql.NullInt32{Int32: 3, Valid: true},
		NoOfBathRooms: sql.NullInt32{Int32: 2, Valid: true},
		NoOfToilets:   sql.NullInt32{Int32: 2, Valid: true},
		GeoLocation:   pqtype.NullRawMessage{RawMessage: []byte(`{"lat": 40.7128, "lng": -74.0060}`), Valid: true},
		Status:        NullPropertyStatus{PropertyStatus: "Available", Valid: true},
	}

	property, err := testQueries.InsertProperty(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, property)

	require.Equal(t, arg.Title, property.Title)
	require.Equal(t, arg.Price, property.Price)
	require.Equal(t, arg.Type, property.Type)
	require.Equal(t, arg.OwnerID, property.OwnerID)
	require.WithinDuration(t, time.Now(), property.CreatedAt.Time, time.Second)

	return property
}

func TestInsertProperty(t *testing.T) {
	property := createRandomProperty(t)
	require.NotEmpty(t, property)
}

func TestGetPropertyByID(t *testing.T) {
	property := createRandomProperty(t)

	retrievedProperty, err := testQueries.GetPropertyByID(context.Background(), property.ID)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedProperty)

	require.Equal(t, property.ID, retrievedProperty.ID)
	require.Equal(t, property.Title, retrievedProperty.Title)
	require.Equal(t, property.Price, retrievedProperty.Price)
	require.WithinDuration(t, property.CreatedAt.Time, retrievedProperty.CreatedAt.Time, time.Second)
}

func TestUpdateProperty(t *testing.T) {
	property := createRandomProperty(t)

	updateArg := UpdatePropertyParams{
		ID:            property.ID,
		Title:         "Updated Property",
		Description:   sql.NullString{String: "Updated description", Valid: true},
		Price:         "600000", // Update to a new valid Price
		Type:          NullPropertyType{PropertyType: "Apartment", Valid: true},
		Address:       "456 Updated St",
		ZipCode:       sql.NullString{String: "54321", Valid: true},
		Images:        []string{"updated_image1.jpg", "updated_image2.jpg"},
		NoOfBedRooms:  sql.NullInt32{Int32: 4, Valid: true},
		NoOfBathRooms: sql.NullInt32{Int32: 3, Valid: true},
		NoOfToilets:   sql.NullInt32{Int32: 3, Valid: true},
		GeoLocation:   pqtype.NullRawMessage{RawMessage: []byte(`{"lat": 37.7749, "lng": -122.4194}`), Valid: true},
		Status:        NullPropertyStatus{PropertyStatus: "Available", Valid: true},
	}

	err := testQueries.UpdateProperty(context.Background(), updateArg)
	require.NoError(t, err)

	updatedProperty, err := testQueries.GetPropertyByID(context.Background(), property.ID)
	require.NoError(t, err)
	require.Equal(t, updateArg.Title, updatedProperty.Title)
	require.Equal(t, updateArg.Price, updatedProperty.Price)
	require.WithinDuration(t, property.CreatedAt.Time, updatedProperty.CreatedAt.Time, time.Second)
}

func TestDeleteProperty(t *testing.T) {
	property := createRandomProperty(t)

	err := testQueries.DeleteProperty(context.Background(), property.ID)
	require.NoError(t, err)

	// Verify the property no longer exists
	_, err = testQueries.GetPropertyByID(context.Background(), property.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListProperties(t *testing.T) {
	// Insert multiple properties
	for i := 0; i < 5; i++ {
		createRandomProperty(t)
	}

	arg := ListPropertiesParams{
		Limit:  5,
		Offset: 0,
	}

	properties, err := testQueries.ListProperties(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, properties, 5)

	for _, property := range properties {
		require.NotEmpty(t, property)
	}
}

func TestGetPropertiesByOwnerID(t *testing.T) {
	ownerID := uuid.New()
	// Create properties with the same owner ID
	for i := 0; i < 3; i++ {
		createRandomPropertyWithOwner(t, ownerID)
	}

	arg := GetPropertiesByOwnerIDParams{
		OwnerID: uuid.NullUUID{UUID: ownerID, Valid: true},
		Limit:   3,
		Offset:  0,
	}

	properties, err := testQueries.GetPropertiesByOwnerID(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, properties, 3)

	for _, property := range properties {
		require.Equal(t, ownerID, property.OwnerID.UUID)
	}
}

// Helper function to create properties with a specific owner ID
func createRandomPropertyWithOwner(t *testing.T, ownerID uuid.UUID) Property {
	arg := InsertPropertyParams{
		ID:            uuid.New(),
		Title:         "Sample Property",
		Description:   sql.NullString{String: "A beautiful property.", Valid: true},
		Price:         "500000",
		Type:          NullPropertyType{PropertyType: "House", Valid: true},
		Address:       "123 Main St",
		ZipCode:       sql.NullString{String: "12345", Valid: true},
		OwnerID:       uuid.NullUUID{UUID: ownerID, Valid: true},
		Images:        []string{"image1.jpg", "image2.jpg"},
		NoOfBedRooms:  sql.NullInt32{Int32: 3, Valid: true},
		NoOfBathRooms: sql.NullInt32{Int32: 2, Valid: true},
		NoOfToilets:   sql.NullInt32{Int32: 2, Valid: true},
		GeoLocation:   pqtype.NullRawMessage{RawMessage: []byte(`{"lat": 40.7128, "lng": -74.0060}`), Valid: true},
		Status:        NullPropertyStatus{PropertyStatus: "Available", Valid: true},
	}

	property, err := testQueries.InsertProperty(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, property)
	return property
}
