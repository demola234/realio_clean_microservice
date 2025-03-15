from abc import ABC, abstractmethod
from typing import List, Optional, Any
from domain.entities.property import Property
from domain.value_objects.user_preference import UserPreference


class PropertyRepository(ABC):
    """Interface for property data repositories"""
    
    @abstractmethod
    def get_all(self) -> List[Property]:
        """Get all properties"""
        pass
    
    @abstractmethod
    def get_by_id(self, property_id: int) -> Optional[Property]:
        """Get a property by ID"""
        pass
    
    @abstractmethod
    def add(self, property_data: Property) -> Property:
        """Add a new property"""
        pass
    
    @abstractmethod
    def update(self, property_data: Property) -> Property:
        """Update an existing property"""
        pass
    
    @abstractmethod
    def delete(self, property_id: int) -> bool:
        """Delete a property by ID"""
        pass
    
    @abstractmethod
    def bulk_update(self, properties: List[Property]) -> int:
        """Update multiple properties at once"""
        pass
    
    @abstractmethod
    def save(self) -> bool:
        """Save changes to the underlying storage"""
        pass
    
    @abstractmethod
    def load(self) -> bool:
        """Load data from the underlying storage"""
        pass


class ModelRepository(ABC):
    """Interface for ML model repositories"""
    
    @abstractmethod
    def save_model(self, model: Any) -> bool:
        """Save a trained model"""
        pass
    
    @abstractmethod
    def load_model(self) -> Any:
        """Load a trained model"""
        pass
    
    @abstractmethod
    def save_encoder(self, encoder: Any) -> bool:
        """Save a label encoder"""
        pass
    
    @abstractmethod
    def load_encoder(self) -> Any:
        """Load a label encoder"""
        pass
    
    @abstractmethod
    def save_scaler(self, scaler: Any) -> bool:
        """Save a feature scaler"""
        pass
    
    @abstractmethod
    def load_scaler(self) -> Any:
        """Load a feature scaler"""
        pass