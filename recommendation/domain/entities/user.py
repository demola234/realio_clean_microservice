from dataclasses import dataclass
from typing import List, Dict, Any


@dataclass
class User:
    """User entity representing a system user"""
    id: int
    favorites: List[int]
    viewed_properties: List[int]
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'User':
        """Create a User entity from a dictionary"""
        return cls(
            id=data.get('user_id', 0),
            favorites=data.get('favorites', []),
            viewed_properties=data.get('viewed_properties', [])
        )
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert User entity to a dictionary"""
        return {
            'user_id': self.id,
            'favorites': self.favorites,
            'viewed_properties': self.viewed_properties
        }