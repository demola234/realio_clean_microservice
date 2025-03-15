from typing import Dict, Any
from dependency_injector import containers, providers

# Import application layer
from application.use_cases.get_recommendations import GetRecommendationsUseCase
from application.use_cases.update_property import UpdatePropertyUseCase
from application.use_cases.retrain_model import RetrainModelUseCase
from application.services.recommendation_service import RecommendationServiceImpl

# Import infrastructure layer
from infrastructure.repositories.file_property_repository import FilePropertyRepository
from infrastructure.repositories.model_repository import FileModelRepository
from infrastructure.ml.preprocessing import PropertyPreprocessor
from infrastructure.message_bus.kafka_message_bus import KafkaMessageBus

# Import interface layer
from interface.controllers.recommendation_controller import RecommendationController
from interface.controllers.admin_controller import AdminController
from interface.event_handlers.property_event_handler import PropertyEventHandler
from interface.event_handlers.model_event_handler import ModelEventHandler
from interface.serializers.grpc_serializer import GrpcSerializer
from interface.api.grpc.recommendation_service import GrpcServer

# Import configuration
from config.settings import Settings


class Container(containers.DeclarativeContainer):
    """Dependency injection container"""
    
    config = providers.Singleton(Settings)
    
    # Infrastructure layer
    property_repository = providers.Singleton(
        FilePropertyRepository,
        file_path=config.provided.property_data_path
    )
    
    model_repository = providers.Singleton(
        FileModelRepository,
        model_path=config.provided.model_path,
        encoder_path=config.provided.encoder_path,
        scaler_path=config.provided.scaler_path
    )
    
    preprocessing_service = providers.Singleton(PropertyPreprocessor)
    
    message_bus = providers.Singleton(
        KafkaMessageBus,
        bootstrap_servers=config.provided.kafka_bootstrap_servers,
        group_id=config.provided.kafka_group_id,
        client_id=config.provided.kafka_client_id
    )
    
    # Application layer
    recommendation_service = providers.Singleton(
        RecommendationServiceImpl,
        property_repository=property_repository,
        model_repository=model_repository,
        preprocessing_service=preprocessing_service
    )
    
    get_recommendations_use_case = providers.Factory(
        GetRecommendationsUseCase,
        property_repository=property_repository,
        recommendation_service=recommendation_service
    )
    
    update_property_use_case = providers.Factory(
        UpdatePropertyUseCase,
        property_repository=property_repository
    )
    
    retrain_model_use_case = providers.Factory(
        RetrainModelUseCase,
        recommendation_service=recommendation_service
    )
    
    # Interface layer
    recommendation_controller = providers.Factory(
        RecommendationController,
        get_recommendations_use_case=get_recommendations_use_case
    )
    
    admin_controller = providers.Factory(
        AdminController,
        retrain_model_use_case=retrain_model_use_case,
        admin_token=config.provided.admin_token
    )
    
    property_event_handler = providers.Factory(
        PropertyEventHandler,
        update_property_use_case=update_property_use_case
    )
    
    model_event_handler = providers.Factory(
        ModelEventHandler,
        retrain_model_use_case=retrain_model_use_case
    )
    
    grpc_serializer = providers.Factory(GrpcSerializer)
    
    grpc_server = providers.Singleton(
        GrpcServer,
        recommendation_controller=recommendation_controller,
        admin_controller=admin_controller,
        serializer=grpc_serializer,
        port=config.provided.grpc_port,
        max_workers=config.provided.grpc_max_workers
    )