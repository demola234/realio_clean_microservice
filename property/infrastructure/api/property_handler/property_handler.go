package propertyhandler

import (
	"context"
	"encoding/json"
	pb "job_portal/property/infrastructure/api/grpc"
	"job_portal/property/internal/domain/entity"
	"job_portal/property/internal/usecases"
	"strconv"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PropertyHandler struct {
	propertyUsecase usecases.PropertyUsecase
	pb.UnimplementedPropertyServiceServer
}

func NewPropertyHandler(propertyUsecase usecases.PropertyUsecase) *PropertyHandler {
	return &PropertyHandler{
		propertyUsecase: propertyUsecase,
	}

}

func (p *PropertyHandler) CreateProperty(ctx context.Context, req *pb.CreatePropertyRequest) (*pb.CreatePropertyResponse, error) {
	userID, err := uuid.Parse(req.GetOwnerId())
	if err != nil {
		return nil, status.Errorf(400, "Invalid user ID")
	}

	// Implement the CreateProperty method
	property := &entity.Property{
		Title:         req.GetTitle(),
		Description:   req.GetDescription(),
		Price:         strconv.FormatFloat(req.GetPrice(), 'f', -1, 64),
		Type:          req.GetType(),
		Address:       req.GetAddress(),
		ZipCode:       req.GetZipCode(),
		Images:        req.GetImages(),
		NoOfBedRooms:  strconv.Itoa(int(req.NoOfBedrooms)),
		NoOfBathRooms: strconv.Itoa(int(req.NoOfBedrooms)),
		NoOfToilets:   strconv.Itoa(int(req.NoOfToilets)),
		GeoLocation: pqtype.NullRawMessage{
			RawMessage: json.RawMessage(req.GetGeoLocation()),
			Valid:      true,
		},
		OwnerID: uuid.NullUUID{UUID: userID, Valid: true},
		Status:  req.GetStatus(),
	}

	p.propertyUsecase.CreateProperty(ctx, property)

	propertyResponse := &pb.CreatePropertyResponse{
		Id: property.ID.String(),
	}

	return propertyResponse, nil
}

func (p *PropertyHandler) GetProperties(ctx context.Context, req *pb.GetPropertiesRequest) (*pb.GetPropertiesResponse, error) {
	// Implement the GetProperties method

	properties, err := p.propertyUsecase.GetProperties(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	propertyResponses := make([]*pb.Property, len(properties))
	for i, property := range properties {
		price, err := strconv.ParseFloat(property.Price, 64)
		if err != nil {
			return nil, status.Errorf(500, "Unable to parse price: %v", err)
		}

		noOfBedRooms, err := strconv.Atoi(property.NoOfBedRooms)
		if err != nil {
			return nil, status.Errorf(500, "Unable to parse no of bedrooms: %v", err)
		}

		noOfBathRooms, err := strconv.Atoi(property.NoOfBathRooms)
		if err != nil {
			return nil, status.Errorf(500, "Unable to parse no of bathrooms: %v", err)
		}

		noOfToilets, err := strconv.Atoi(property.NoOfToilets)
		if err != nil {
			return nil, status.Errorf(500, "Unable to parse no of toilets: %v", err)
		}

		propertyResponses[i] = &pb.Property{
			Id:            property.ID.String(),
			Title:         property.Title,
			Description:   property.Description,
			Price:         price,
			Type:          property.Type,
			Address:       property.Address,
			ZipCode:       property.ZipCode,
			Images:        property.Images,
			OwnerId:       property.OwnerID.UUID.String(),
			NoOfBedrooms:  int32(noOfBedRooms),
			NoOfBathrooms: int32(noOfBathRooms),
			NoOfToilets:   int32(noOfToilets),
			GeoLocation:   string(property.GeoLocation.RawMessage),
			Status:        property.Status,
			CreatedAt:     timestamppb.New(property.CreatedAt),
			UpdatedAt:     timestamppb.New(property.UpdatedAt),
		}
	}

	return &pb.GetPropertiesResponse{
		Properties: propertyResponses,
	}, nil
}

func (p *PropertyHandler) GetPropertiesByOwner(ctx context.Context, req *pb.GetPropertiesByOwnerRequest) (*pb.GetPropertiesByOwnerResponse, error) {
	// Implement the GetProperties method
	// Parse req.OwnerId() to uuuid.UUID
	uuidUser, err := uuid.Parse(req.GetOwnerId())
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	properties, err := p.propertyUsecase.GetPropertiesByOwner(ctx, uuid.NullUUID{UUID: uuidUser, Valid: true}, req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	propertyResponses := make([]*pb.Property, len(properties))
	for i, property := range properties {
		price, err := strconv.ParseFloat(property.Price, 64)
		if err != nil {
			return nil, status.Errorf(500, "Unable to parse price: %v", err)
		}

		noOfBedRooms, err := strconv.Atoi(property.NoOfBedRooms)
		if err != nil {
			return nil, status.Errorf(500, "Unable to parse no of bedrooms: %v", err)
		}

		noOfBathRooms, err := strconv.Atoi(property.NoOfBathRooms)
		if err != nil {
			return nil, status.Errorf(500, "Unable to parse no of bathrooms: %v", err)
		}

		noOfToilets, err := strconv.Atoi(property.NoOfToilets)
		if err != nil {
			return nil, status.Errorf(500, "Unable to parse no of toilets: %v", err)
		}

		propertyResponses[i] = &pb.Property{
			Id:            property.ID.String(),
			Title:         property.Title,
			Description:   property.Description,
			Price:         price,
			Type:          property.Type,
			Address:       property.Address,
			ZipCode:       property.ZipCode,
			Images:        property.Images,
			OwnerId:       property.OwnerID.UUID.String(),
			NoOfBedrooms:  int32(noOfBedRooms),
			NoOfBathrooms: int32(noOfBathRooms),
			NoOfToilets:   int32(noOfToilets),
			GeoLocation:   string(property.GeoLocation.RawMessage),
			Status:        property.Status,
			CreatedAt:     timestamppb.New(property.CreatedAt),
			UpdatedAt:     timestamppb.New(property.UpdatedAt),
		}
	}

	return &pb.GetPropertiesByOwnerResponse{
		Properties: propertyResponses,
	}, nil
}

func (p *PropertyHandler) GetPropertyByID(ctx context.Context, req *pb.GetPropertyByIDRequest) (*pb.GetPropertyByIDResponse, error) {
	// Implement the GetPropertyByID method
	uuidProperty, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}
	property, err := p.propertyUsecase.GetPropertyByID(ctx, uuidProperty)
	if err != nil {
		return nil, status.Errorf(500, "Unable to get property: %v", err)
	}

	price, err := strconv.ParseFloat(property.Price, 64)
	if err != nil {
		return nil, status.Errorf(500, "Unable to parse price: %v", err)
	}

	noOfBedRooms, err := strconv.Atoi(property.NoOfBedRooms)
	if err != nil {
		return nil, status.Errorf(500, "Unable to parse no of bedrooms: %v", err)
	}

	noOfBathRooms, err := strconv.Atoi(property.NoOfBathRooms)
	if err != nil {
		return nil, status.Errorf(500, "Unable to parse no of bathrooms: %v", err)
	}

	noOfToilets, err := strconv.Atoi(property.NoOfToilets)
	if err != nil {
		return nil, status.Errorf(500, "Unable to parse no of toilets: %v", err)
	}

	propertyResponse := &pb.Property{
		Id:            property.ID.String(),
		Title:         property.Title,
		Description:   property.Description,
		Price:         price,
		Type:          property.Type,
		Address:       property.Address,
		ZipCode:       property.ZipCode,
		Images:        property.Images,
		OwnerId:       property.OwnerID.UUID.String(),
		NoOfBedrooms:  int32(noOfBedRooms),
		NoOfBathrooms: int32(noOfBathRooms),
		NoOfToilets:   int32(noOfToilets),
		GeoLocation:   string(property.GeoLocation.RawMessage),
		Status:        property.Status,
		CreatedAt:     timestamppb.New(property.CreatedAt),
		UpdatedAt:     timestamppb.New(property.UpdatedAt),
	}

	return &pb.GetPropertyByIDResponse{
		Property: propertyResponse,
	}, nil

}

func (p *PropertyHandler) UpdateProperty(ctx context.Context, req *pb.UpdatePropertyRequest) (*pb.UpdatePropertyResponse, error) {
	// Implement the UpdateProperty method
	uuidProperty, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	// Get property by id
	UserProperty, err := p.propertyUsecase.GetPropertyByID(ctx, uuidProperty)
	if err != nil {
		return nil, status.Errorf(500, "Unable to get property: %v", err)
	}

	// Allow only owner to update property
	if UserProperty.OwnerID.UUID.String() != req.GetOwnerId() {
		return nil, status.Errorf(401, "You are not authorized to update this property")
	}

	property := &entity.Property{
		ID:            uuidProperty,
		Title:         req.GetTitle(),
		Description:   req.GetDescription(),
		Price:         strconv.FormatFloat(req.GetPrice(), 'f', -1, 64),
		Type:          req.GetType(),
		Address:       req.GetAddress(),
		ZipCode:       req.GetZipCode(),
		Images:        req.GetImages(),
		NoOfBedRooms:  strconv.Itoa(int(req.NoOfBedrooms)),
		NoOfBathRooms: strconv.Itoa(int(req.NoOfBedrooms)),
		NoOfToilets:   strconv.Itoa(int(req.NoOfToilets)),
		GeoLocation: pqtype.NullRawMessage{
			RawMessage: json.RawMessage(req.GetGeoLocation()),
			Valid:      true,
		},
		Status: req.GetStatus(),
	}

	err = p.propertyUsecase.UpdateProperty(ctx, property)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePropertyResponse{}, nil

}

func (p *PropertyHandler) DeleteProperty(ctx context.Context, req *pb.DeletePropertyRequest) (*pb.DeletePropertyResponse, error) {
	// Implement the DeleteProperty method
	uuidProperty, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	// Get property by id
	UserProperty, err := p.propertyUsecase.GetPropertyByID(ctx, uuidProperty)
	if err != nil {
		return nil, status.Errorf(500, "Unable to get property: %v", err)
	}

	// Allow only owner to update property
	if UserProperty.OwnerID.UUID.String() != req.GetOwnerId() {
		return nil, status.Errorf(401, "You are not authorized to update this property")
	}

	err = p.propertyUsecase.DeleteProperty(ctx, uuidProperty)

	if err != nil {
		return nil, err
	}

	return &pb.DeletePropertyResponse{}, nil
}
