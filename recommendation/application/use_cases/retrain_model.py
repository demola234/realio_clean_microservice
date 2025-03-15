
# application/use_cases/retrain_model.py
from application.interfaces.services import RecommendationService


class RetrainModelUseCase:
    """Retrain the recommendation model"""
    
    def __init__(self, recommendation_service: RecommendationService):
        self.recommendation_service = recommendation_service
    
    def execute(self, force: bool = False) -> bool:
        """Execute the use case
        
        Args:
            force: Force retraining even if not necessary
            
        Returns:
            True if model was retrained, False otherwise
        """
        # Retrain the model
        success = self.recommendation_service.retrain_model(force)
        
        return success

