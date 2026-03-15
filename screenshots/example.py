from dataclasses import dataclass, field
from typing import (
    Any,
    Dict,
)
from uuid import UUID, uuid4


@dataclass
class Example:
    """An example class."""

    name: str = field(default="John")
    """The name of the user"""

    age: int = field(
        default=20, init=False, metadata={"description": "The age of the user"}
    )
    """The age of the user."""

    id: UUID = field(default_factory=uuid4)
    """The ID of the user."""

    profile: Dict[str, Any] | None = None
    """The profile of the user."""

    def get_user(self) -> Dict[str, Any]:
        return {"name": self.name, "age": self.age, "id": self.id}
