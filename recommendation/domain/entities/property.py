from dataclasses import dataclass
from typing import Optional, Dict, Any


@dataclass
class Property:
    """Property entity representing a real estate property"""
    id: int
    location: str
    bedrooms: int
    bathrooms: int
    toilets: int
    parking_spaces: Optional[int] = None
    price: Optional[float] = None
    predicted_price: Optional[float] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'Property':
        """Create a Property entity from a dictionary"""
        return cls(
            id=data.get('id', 0),
            location=data.get('location', ''),
            bedrooms=data.get('Bedrooms', 0),
            bathrooms=data.get('Bathrooms', 0),
            toilets=data.get('Toilets', 0),
            parking_spaces=data.get('Parking Spaces'),
            price=data.get('price'),
            predicted_price=data.get('predicted_price')
        )
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert Property entity to a dictionary"""
        return {
            'id': self.id,
            'location': self.location,
            'Bedrooms': self.bedrooms,
            'Bathrooms': self.bathrooms,
            'Toilets': self.toilets,
            'Parking Spaces': self.parking_spaces,
            'price': self.price,
            'predicted_price': self.predicted_price
        }

