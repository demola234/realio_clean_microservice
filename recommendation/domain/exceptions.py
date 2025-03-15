class DomainException(Exception):
    """Base exception for all domain exceptions"""
    pass


class PropertyNotFoundException(DomainException):
    """Exception raised when a property is not found"""
    pass


class UserNotFoundException(DomainException):
    """Exception raised when a user is not found"""
    pass


class InvalidPropertyDataException(DomainException):
    """Exception raised when property data is invalid"""
    pass


class ModelTrainingException(DomainException):
    """Exception raised when model training fails"""
    pass
