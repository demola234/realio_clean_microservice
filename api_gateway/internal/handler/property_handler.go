package handler

import (
	"context"
	"job_portal/api_gateway/interfaces/grpc_clients"
	token "job_portal/api_gateway/interfaces/middleware/token_maker"
	pb "job_portal/property/interfaces/api/grpc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PropertyHandler struct {
	PropertyClient *grpc_clients.PropertyClient
}

func NewPropertyHandler(propertyClient *grpc_clients.PropertyClient) *PropertyHandler {
	return &PropertyHandler{PropertyClient: propertyClient}
}

func (h *PropertyHandler) GetProperties(c *gin.Context) {

	// Get limit and offset from query parameters and convert them to int32
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit value"})
		return
	}
	limitInt32 := int32(limit)

	offsetStr := c.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset value"})
		return
	}
	offsetInt32 := int32(offset)

	// Call the gRPC client with the converted parameters
	res, err := h.PropertyClient.Client.GetProperties(context.Background(), &pb.GetPropertiesRequest{
		Limit:  limitInt32,
		Offset: offsetInt32,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *PropertyHandler) GetPropertiesByOwner(c *gin.Context) {
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	userID := authPayload.(*token.Payload).UserID

	// Get limit and offset from query parameters and convert them to int32
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit value"})
		return
	}
	limitInt32 := int32(limit)

	offsetStr := c.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset value"})
		return
	}
	offsetInt32 := int32(offset)

	// Call the gRPC client with the converted parameters
	res, err := h.PropertyClient.Client.GetPropertiesByOwner(context.Background(), &pb.GetPropertiesByOwnerRequest{
		OwnerId: userID,
		Limit:   limitInt32,
		Offset:  offsetInt32,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *PropertyHandler) GetProperty(c *gin.Context) {
	propertyID := c.Param("id")

	res, err := h.PropertyClient.Client.GetPropertyByID(context.Background(), &pb.GetPropertyByIDRequest{Id: propertyID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	userID := authPayload.(*token.Payload).UserID

	var req pb.CreatePropertyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.OwnerId = userID

	res, err := h.PropertyClient.Client.CreateProperty(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	userID := authPayload.(*token.Payload).UserID
	propertyID := c.Param("id")
	var req pb.UpdatePropertyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.OwnerId = userID
	req.Id = propertyID
	res, err := h.PropertyClient.Client.UpdateProperty(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
