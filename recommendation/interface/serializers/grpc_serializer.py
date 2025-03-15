# interface/serializers/grpc_serializer.py
from typing import Dict, Any, List
from interface.api.grpc import recommendation_pb2


class GrpcSerializer:
    """Serializer for gRPC messages"""
    
    def serialize_recommendation_response(self, data: Dict[str, Any]) -> recommendation_pb2.RecommendationResponse:
        """Serialize recommendation data to a gRPC RecommendationResponse
        
        Args:
            data: Dictionary with status, message, and recommendations
            
        Returns:
            RecommendationResponse gRPC message
        """
        # Create response object
        response = recommendation_pb2.RecommendationResponse(
            status=data['status'],
            message=data['message']
        )
        
        # Add recommendations
        for rec in data['recommendations']:
            property_msg = recommendation_pb2.Property(
                property_id=rec['property_id'],
                location=rec['location'],
                bedrooms=rec['bedrooms'],
                bathrooms=rec['bathrooms'],
                toilets=rec['toilets'],
                price=rec['price'],
                predicted_price=rec['predicted_price']
            )
            response.recommendations.append(property_msg)
        
        return response