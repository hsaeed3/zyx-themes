from dataclasses import dataclass, field
import functools
from typing import (
    Any,
    Dict,
    TypeVar,
)
from uuid import UUID, uuid4


_T = TypeVar("_T")


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

    @property
    def full_name(self) -> str:
        return f"{self.name} {self.age}"

    def get_user(self) -> Dict[str, Any]:
        return {"name": self.name, "age": self.age, "id": self.id}

    @functools.lru_cache()
    async def async_function(self) -> int:
        return 2