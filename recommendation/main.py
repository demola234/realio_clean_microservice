from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import pandas as pd
import joblib
import numpy as np
import re
from sklearn.preprocessing import StandardScaler

# Initialize FastAPI app
app = FastAPI()

# Load model and preprocessing objects
try:
    recommendation_model = joblib.load("recommendation_model.pkl")
    label_encoder = joblib.load("label_encoder.pkl")
    scaler = joblib.load("scaler.pkl")
except Exception as e:
    print(f"Error loading models: {e}")

# Load and preprocess property data
try:
    property_data = pd.read_json('properties.json')
    
    def clean_rooms_info(rooms):
        if pd.isna(rooms) or not isinstance(rooms, str):
            return [np.nan, np.nan, np.nan, np.nan]
        
        rooms = rooms.replace(' Save', '').strip()
        bedrooms = bathrooms = toilets = parking_spaces = np.nan

        bedrooms_match = re.search(r'(\d+)\s*Bedrooms?', rooms)
        bathrooms_match = re.search(r'(\d+)\s*Bathrooms?', rooms)
        toilets_match = re.search(r'(\d+)\s*Toilets?', rooms)
        parking_spaces_match = re.search(r'(\d+)\s*Parking Spaces?', rooms)

        if bedrooms_match:
            bedrooms = int(bedrooms_match.group(1))
        if bathrooms_match:
            bathrooms = int(bathrooms_match.group(1))
        if toilets_match:
            toilets = int(toilets_match.group(1))
        if parking_spaces_match:
            parking_spaces = int(parking_spaces_match.group(1))

        return [bedrooms, bathrooms, toilets, parking_spaces]

    property_data[['Bedrooms', 'Bathrooms', 'Toilets', 'Parking Spaces']] = property_data['rooms'].apply(
        lambda x: pd.Series(clean_rooms_info(x))
    )

    for column in ['Bedrooms', 'Bathrooms', 'Toilets', 'Parking Spaces']:
        property_data[column] = property_data[column].fillna(property_data[column].median())
    
    property_data['price'] = property_data['price'].replace('[\₦\$,]', '', regex=True).astype(float)
    property_data['location_encoded'] = label_encoder.fit_transform(property_data['location'])

    # Define X_clean as the features used for model predictions
    X_clean = property_data[['Bedrooms', 'Bathrooms', 'Toilets', 'Parking Spaces', 'location_encoded']]

except Exception as e:
    print(f"Error loading or processing property data: {e}")

# Define input schema
class UserData(BaseModel):
    user_id: int
    favorites: list[int]
    viewed_properties: list[int]

def clean_user_preferences(user_data, property_data):
    behavior_set = set(user_data.favorites + user_data.viewed_properties)
    behavior_properties = property_data[property_data.index.isin(behavior_set)]
    
    user_profile = {
        'Bedrooms': behavior_properties['Bedrooms'].median(),
        'Bathrooms': behavior_properties['Bathrooms'].median(),
        'Toilets': behavior_properties['Toilets'].median(),
        'Parking Spaces': behavior_properties['Parking Spaces'].median(),
        'location': behavior_properties['location'].mode()[0]
    } if not behavior_properties.empty else {
        'Bedrooms': property_data['Bedrooms'].median(),
        'Bathrooms': property_data['Bathrooms'].median(),
        'Toilets': property_data['Toilets'].median(),
        'Parking Spaces': property_data['Parking Spaces'].median(),
        'location': property_data['location'].mode()[0]
    }
    
    return user_profile

def recommend_properties(user_preferences, property_data, n_recommendations=5):
    try:
        closest_location = label_encoder.transform([user_preferences['location']])[0]
    except:
        raise HTTPException(status_code=404, detail="Location not found in dataset.")

    example_features = pd.DataFrame([{
        "Bedrooms": user_preferences['Bedrooms'],
        "Bathrooms": user_preferences['Bathrooms'],
        "Toilets": user_preferences['Toilets'],
        "Parking Spaces": user_preferences['Parking Spaces'],
        "location_encoded": closest_location
    }])

    # Scale the features for price prediction
    example_features_scaled = scaler.transform(example_features[['Bedrooms', 'Bathrooms', 'Toilets', 'Parking Spaces']])
    example_features_scaled = np.concatenate([example_features_scaled, example_features[['location_encoded']].values], axis=1)

    # Predict prices using the recommendation model
    property_data['predicted_price'] = recommendation_model.predict(X_clean)

    # Filter for unique locations and user preferences
    recommended_properties = property_data[
        (property_data['Bedrooms'] == example_features.iloc[0]['Bedrooms']) &
        (property_data['Bathrooms'] == example_features.iloc[0]['Bathrooms']) &
        (property_data['Toilets'] == example_features.iloc[0]['Toilets']) &
        (property_data['location_encoded'] == closest_location)
    ].drop_duplicates(subset=['location'])

    # Sort by predicted price, add random properties if necessary
    recommended_properties = recommended_properties.sort_values(by='predicted_price').head(n_recommendations)
    
    if len(recommended_properties) < n_recommendations:
        remaining_count = n_recommendations - len(recommended_properties)
        additional_properties = property_data[~property_data['location'].isin(recommended_properties['location'])].sample(remaining_count)
        recommended_properties = pd.concat([recommended_properties, additional_properties])

    return recommended_properties[['location', 'Bedrooms', 'Bathrooms', 'Toilets', 'price', 'predicted_price']].head(n_recommendations)

# API endpoint
@app.post("/recommend")
async def get_recommendations(user_data: UserData):
    try:
        user_preferences = clean_user_preferences(user_data, property_data)
        recommendations = recommend_properties(user_preferences, property_data)
        return recommendations.to_dict(orient="records")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
