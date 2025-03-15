from abc import ABC, abstractmethod
from typing import List, Dict, Any
from domain.entities.property import Property
from domain.entities.user import User
from domain.value_objects.user_preference import UserPreference


class RecommendationService(ABC):
    """Interface for recommendation services"""
    
    @abstractmethod
    def get_recommendations(self, user: User, limit: int = 5) -> List[Property]:
        """Get property recommendations for a user"""
        pass
    
    @abstractmethod
    def extract_user_preferences(self, user: User, properties: List[Property]) -> UserPreference:
        """Extract user preferences from their behavior"""
        pass
    
    @abstractmethod
    def retrain_model(self, force: bool = False) -> bool:
        """Retrain the recommendation model"""
        pass
    
    @abstractmethod
    def predict_prices(self, properties: List[Property]) -> List[float]:
        """Predict prices for a list of properties"""
        pass


class PreprocessingService(ABC):
    """Interface for data preprocessing services"""
    
    @abstractmethod
    def preprocess_properties(self, properties: List[Property], for_training: bool = False) -> Dict[str, Any]:
        """Preprocess property data for model training or prediction"""
        pass

