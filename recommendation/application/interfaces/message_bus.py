
# application/interfaces/message_bus.py
from abc import ABC, abstractmethod
from typing import Callable, Dict, Any


class MessagePublisher(ABC):
    """Interface for publishing messages"""
    
    @abstractmethod
    def publish(self, topic: str, message: Dict[str, Any]) -> bool:
        """Publish a message to a topic"""
        pass


class MessageConsumer(ABC):
    """Interface for consuming messages"""
    
    @abstractmethod
    def subscribe(self, topic: str, handler: Callable[[Dict[str, Any]], None]) -> None:
        """Subscribe to a topic with a message handler"""
        pass
    
    @abstractmethod
    def start(self) -> None:
        """Start consuming messages"""
        pass
    
    @abstractmethod
    def stop(self) -> None:
        """Stop consuming messages"""
        pass


class MessageBus(MessagePublisher, MessageConsumer):
    """Combined interface for message publishing and consuming"""
    pass