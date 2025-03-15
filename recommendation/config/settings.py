import os
from typing import Dict, Any


class Settings:
    """Application settings"""
    
    def __init__(self):
        # Property repository settings
        self.property_data_path = os.environ.get('PROPERTY_DATA_PATH', 'data/properties.json')
        
        # Model repository settings
        self.model_path = os.environ.get('MODEL_PATH', 'models/recommendation_model.pkl')
        self.encoder_path = os.environ.get('ENCODER_PATH', 'models/label_encoder.pkl')
        self.scaler_path = os.environ.get('SCALER_PATH', 'models/scaler.pkl')
        
        # Kafka settings
        self.kafka_bootstrap_servers = os.environ.get('KAFKA_BOOTSTRAP_SERVERS', 'localhost:9092')
        self.kafka_group_id = os.environ.get('KAFKA_GROUP_ID', 'recommendation_service')
        self.kafka_client_id = os.environ.get('KAFKA_CLIENT_ID', 'recommendation_service')
        self.kafka_property_topic = os.environ.get('KAFKA_PROPERTY_TOPIC', 'property_updates')
        self.kafka_model_topic = os.environ.get('KAFKA_MODEL_TOPIC', 'model_retrain')
        
        # gRPC settings
        self.grpc_port = int(os.environ.get('GRPC_PORT', '50051'))
        self.grpc_max_workers = int(os.environ.get('GRPC_MAX_WORKERS', '10'))
        
        # Admin settings
        self.admin_token = os.environ.get('ADMIN_TOKEN', 'admin-token-123')
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert settings to a dictionary"""
        return {
            'property_data_path': self.property_data_path,
            'model_path': self.model_path,
            'encoder_path': self.encoder_path,
            'scaler_path': self.scaler_path,
            'kafka_bootstrap_servers': self.kafka_bootstrap_servers,
            'kafka_group_id': self.kafka_group_id,
            'kafka_client_id': self.kafka_client_id,
            'kafka_property_topic': self.kafka_property_topic,
            'kafka_model_topic': self.kafka_model_topic,
            'grpc_port': self.grpc_port,
            'grpc_max_workers': self.grpc_max_workers,
            'admin_token': self.admin_token
        }