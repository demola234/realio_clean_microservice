
# application/services/recommendation_service.py
from typing import List, Dict, Any, Optional
import numpy as np
from domain.entities.user import User
from domain.entities.property import Property
from domain.value_objects.user_preference import UserPreference
from application.interfaces.repositories import PropertyRepository, ModelRepository
from application.interfaces.services import RecommendationService, PreprocessingService


class RecommendationServiceImpl(RecommendationService):
    """Implementation of the RecommendationService interface"""
    
    def __init__(
        self, 
        property_repository: PropertyRepository,
        model_repository: ModelRepository,
        preprocessing_service: PreprocessingService
    ):
        self.property_repository = property_repository
        self.model_repository = model_repository
        self.preprocessing_service = preprocessing_service
    
    def get_recommendations(self, user: User, limit: int = 5) -> List[Property]:
        """Get property recommendations for a user"""
        # Get all properties
        all_properties = self.property_repository.get_all()
        
        # Extract user preferences
        user_preferences = self.extract_user_preferences(user, all_properties)
        
        # Preprocess properties for prediction
        processed_data = self.preprocessing_service.preprocess_properties(all_properties)
        
        # Load the model
        model = self.model_repository.load_model()
        scaler = self.model_repository.load_scaler()
        label_encoder = self.model_repository.load_encoder()
        
        if not model or not scaler or not label_encoder:
            # If model is not available, return some properties based on user preferences
            return self._fallback_recommendations(user_preferences, all_properties, limit)
        
        # Predict prices for all properties
        processed_data['properties_df']['predicted_price'] = model.predict(processed_data['X_features'])
        
        # Filter properties based on user preferences
        filtered_properties = self._filter_by_preferences(
            processed_data['properties_df'],
            user_preferences,
            label_encoder
        )
        
        # Sort by predicted price
        if 'predicted_price' in filtered_properties.columns:
            filtered_properties = filtered_properties.sort_values('predicted_price')
        
        # Convert to Property entities
        recommended_properties = []
        for _, row in filtered_properties.iterrows():
            property_dict = row.to_dict()
            property_dict['id'] = row.name if hasattr(row, 'name') else 0
            recommended_properties.append(Property.from_dict(property_dict))
        
        # Limit results
        return recommended_properties[:limit]
    
    def extract_user_preferences(self, user: User, properties: List[Property]) -> UserPreference:
        """Extract user preferences from their behavior"""
        # Get user's favorite and viewed properties
        user_property_ids = set(user.favorites + user.viewed_properties)
        user_properties = [p for p in properties if p.id in user_property_ids]
        
        if not user_properties:
            # If no user properties found, use defaults
            all_bedrooms = [p.bedrooms for p in properties if p.bedrooms is not None]
            all_bathrooms = [p.bathrooms for p in properties if p.bathrooms is not None]
            all_toilets = [p.toilets for p in properties if p.toilets is not None]
            all_parking = [p.parking_spaces for p in properties if p.parking_spaces is not None]
            
            median_bedrooms = np.median(all_bedrooms) if all_bedrooms else 3
            median_bathrooms = np.median(all_bathrooms) if all_bathrooms else 2
            median_toilets = np.median(all_toilets) if all_toilets else 1
            median_parking = np.median(all_parking) if all_parking else 1
            
            # Get most common location
            locations = [p.location for p in properties if p.location]
            most_common_location = max(set(locations), key=locations.count) if locations else "Unknown"
            
            return UserPreference(
                bedrooms=median_bedrooms,
                bathrooms=median_bathrooms,
                toilets=median_toilets,
                parking_spaces=median_parking,
                location=most_common_location
            )
        
        # Calculate median values for numerical preferences
        bedrooms = np.median([p.bedrooms for p in user_properties if p.bedrooms is not None])
        bathrooms = np.median([p.bathrooms for p in user_properties if p.bathrooms is not None])
        toilets = np.median([p.toilets for p in user_properties if p.toilets is not None])
        parking_spaces = np.median([p.parking_spaces for p in user_properties if p.parking_spaces is not None])
        
        # Get most common location
        locations = [p.location for p in user_properties if p.location]
        most_common_location = max(set(locations), key=locations.count) if locations else "Unknown"
        
        return UserPreference(
            bedrooms=bedrooms,
            bathrooms=bathrooms,
            toilets=toilets,
            parking_spaces=parking_spaces,
            location=most_common_location
        )
    
    def retrain_model(self, force: bool = False) -> bool:
        """Retrain the recommendation model"""
        try:
            # Get all properties
            all_properties = self.property_repository.get_all()
            
            # Preprocess properties for training
            processed_data = self.preprocessing_service.preprocess_properties(all_properties, for_training=True)
            
            # Train a new model using scikit-learn
            from sklearn.ensemble import RandomForestRegressor
            model = RandomForestRegressor(n_estimators=100, random_state=42)
            
            # Get training data
            X_train = processed_data['X_features']
            y_train = processed_data['properties_df']['price'].values
            
            # Train the model
            model.fit(X_train, y_train)
            
            # Save the model, encoder, and scaler
            self.model_repository.save_model(model)
            self.model_repository.save_encoder(processed_data['label_encoder'])
            self.model_repository.save_scaler(processed_data['scaler'])
            
            return True
        except Exception as e:
            # Log the error
            print(f"Error training model: {e}")
            return False
    
    def predict_prices(self, properties: List[Property]) -> List[float]:
        """Predict prices for a list of properties"""
        # Preprocess properties
        processed_data = self.preprocessing_service.preprocess_properties(properties)
        
        # Load the model
        model = self.model_repository.load_model()
        
        if not model:
            # Return zeros if model is not available
            return [0.0] * len(properties)
        
        # Predict prices
        predictions = model.predict(processed_data['X_features'])
        
        return predictions.tolist()
    
    def _filter_by_preferences(self, properties_df, user_preferences, label_encoder):
        """Filter properties based on user preferences"""
        import pandas as pd
        
        # Convert location to encoded value
        try:
            location_encoded = label_encoder.transform([user_preferences.location])[0]
        except:
            # If location not found, don't filter by location
            location_encoded = None
        
        # Filter properties
        filtered = properties_df[
            (properties_df['Bedrooms'] == user_preferences.bedrooms) &
            (properties_df['Bathrooms'] == user_preferences.bathrooms) &
            (properties_df['Toilets'] == user_preferences.toilets)
        ]
        
        # Filter by location if available
        if location_encoded is not None and 'location_encoded' in properties_df.columns:
            location_filtered = filtered[filtered['location_encoded'] == location_encoded]
            # If no properties match the location, fall back to the previous filter
            if not location_filtered.empty:
                filtered = location_filtered
        
        # If no properties match the filters, return some random properties
        if filtered.empty and not properties_df.empty:
            return properties_df.sample(min(5, len(properties_df)))
        
        return filtered
    
    def _fallback_recommendations(self, user_preferences, all_properties, limit):
        """Fallback recommendations when model is not available"""
        # Filter properties based on bedrooms and bathrooms
        filtered_properties = [
            p for p in all_properties 
            if p.bedrooms == user_preferences.bedrooms and p.bathrooms == user_preferences.bathrooms
        ]
        
        # If no matches, just return random properties
        if not filtered_properties and all_properties:
            import random
            return random.sample(all_properties, min(limit, len(all_properties)))
        
        # Further filter by location if possible
        location_filtered = [p for p in filtered_properties if p.location == user_preferences.location]
        
        if location_filtered:
            return location_filtered[:limit]
        
        return filtered_properties[:limit]