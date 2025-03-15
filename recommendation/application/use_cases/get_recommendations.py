# application/use_cases/get_recommendations.py
from typing import List, Dict, Any
from domain.entities.user import User
from domain.entities.property import Property
from application.interfaces.repositories import PropertyRepository
from application.interfaces.services import RecommendationService


class GetRecommendationsUseCase:
    """Get property recommendations for a user"""
    
    def __init__(self, property_repository: PropertyRepository, recommendation_service: RecommendationService):
        self.property_repository = property_repository
        self.recommendation_service = recommendation_service
    
    def execute(self, user_data: Dict[str, Any], limit: int = 5) -> List[Property]:
        """Execute the use case
        
        Args:
            user_data: Dictionary containing user data (user_id, favorites, viewed_properties)
            limit: Maximum number of recommendations to return
            
        Returns:
            List of recommended properties
        """
        # Create user entity
        user = User.from_dict(user_data)
        
        # Get recommendations
        recommendations = self.recommendation_service.get_recommendations(user, limit)
        
        return recommendations


# application/use_cases/update_property.py
from typing import List, Dict, Any
from domain.entities.property import Property
from application.interfaces.repositories import PropertyRepository


class UpdatePropertyUseCase:
    """Update property data"""
    
    def __init__(self, property_repository: PropertyRepository):
        self.property_repository = property_repository
    
    def execute(self, property_data: List[Dict[str, Any]]) -> int:
        """Execute the use case
        
        Args:
            property_data: List of property data dictionaries
            
        Returns:
            Number of properties updated
        """
        # Convert dictionaries to Property entities
        properties = [Property.from_dict(data) for data in property_data]
        
        # Update properties in repository
        updated_count = self.property_repository.bulk_update(properties)
        
        # Save changes
        self.property_repository.save()
        
        return updated_count

