from dataclasses import dataclass
from typing import Dict, Any


@dataclass(frozen=True)
class UserPreference:
    """Value object representing a user's property preferences"""
    bedrooms: float
    bathrooms: float
    toilets: float
    parking_spaces: float
    location: str
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert UserPreference to a dictionary"""
        return {
            'Bedrooms': self.bedrooms,
            'Bathrooms': self.bathrooms,
            'Toilets': self.toilets,
            'Parking Spaces': self.parking_spaces,
            'location': self.location
        }
