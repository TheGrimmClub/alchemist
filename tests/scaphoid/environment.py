"""Environment — immutable environment variable set for subprocess calls."""

from __future__ import annotations

import os
from typing import Mapping

__all__ = ["Environment"]


class Environment:
    """An immutable set of environment variables for a subprocess call.

    Two construction modes:

    * :meth:`extended` — starts from ``os.environ`` and adds/overrides keys.
      This is the most common mode: inherit the system environment and inject
      a few test-specific variables on top.
    * :meth:`isolated` — uses *only* the given vars; nothing from
      ``os.environ`` bleeds through, which is useful for hermetic tests.

    Both modes are immutable. :meth:`with_var` always returns a new instance::

        env = Environment.extended({"DEBUG": "1"})
        env2 = env.with_var("PORT", "8080")   # env is unchanged
    """

    __slots__ = ("_mapping", "_mode")

    def __init__(self, mapping: dict[str, str], *, mode: str) -> None:
        object.__setattr__(self, "_mapping", dict(mapping))
        object.__setattr__(self, "_mode", mode)

    def __setattr__(self, name: str, value: object) -> None:
        raise AttributeError("Environment is immutable")

    # --- factories -----------------------------------------------------------

    @classmethod
    def extended(cls, extras: Mapping[str, str] | None = None) -> "Environment":
        """os.environ merged with *extras*; extras win on conflict."""
        return cls(dict(extras or {}), mode="extended")

    @classmethod
    def isolated(cls, vars: Mapping[str, str]) -> "Environment":
        """Use *only* these vars — nothing from os.environ."""
        return cls(dict(vars), mode="isolated")

    # --- mutation (immutable update) -----------------------------------------

    def with_var(self, key: str, value: str) -> "Environment":
        """Return a new Environment with *key* set to *value*."""
        return Environment({**self._mapping, key: value}, mode=self._mode)

    # --- resolution ----------------------------------------------------------

    def resolve(self) -> dict[str, str]:
        """Return the concrete dict to pass to subprocess."""
        if self._mode == "isolated":
            return dict(self._mapping)
        return {**os.environ, **self._mapping}

    # --- introspection -------------------------------------------------------

    @property
    def vars(self) -> dict[str, str]:
        """The explicitly set variables (not including os.environ)."""
        return dict(self._mapping)

    @property
    def is_isolated(self) -> bool:
        return self._mode == "isolated"

    # --- dunder --------------------------------------------------------------

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, Environment):
            return NotImplemented
        return self._mapping == other._mapping and self._mode == other._mode

    def __hash__(self) -> int:
        return hash((tuple(sorted(self._mapping.items())), self._mode))

    def __repr__(self) -> str:
        return f"Environment.{self._mode}({dict(self._mapping)!r})"
