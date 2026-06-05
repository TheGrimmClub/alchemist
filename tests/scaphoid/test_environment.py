import os
import pytest
from scaphoid.environment import Environment


# ---------------------------------------------------------------------------
# extended — construction
# ---------------------------------------------------------------------------

class TestExtended:
    def test_no_args_creates_empty_extras(self):
        env = Environment.extended()
        assert env.vars == {}

    def test_extras_stored(self):
        env = Environment.extended({"KEY": "val"})
        assert env.vars == {"KEY": "val"}

    def test_is_not_isolated(self):
        assert Environment.extended().is_isolated is False

    def test_none_treated_same_as_empty(self):
        assert Environment.extended(None).vars == {}


# ---------------------------------------------------------------------------
# isolated — construction
# ---------------------------------------------------------------------------

class TestIsolated:
    def test_vars_stored(self):
        env = Environment.isolated({"A": "1", "B": "2"})
        assert env.vars == {"A": "1", "B": "2"}

    def test_is_isolated(self):
        assert Environment.isolated({"A": "1"}).is_isolated is True

    def test_empty_vars(self):
        env = Environment.isolated({})
        assert env.vars == {}


# ---------------------------------------------------------------------------
# with_var — immutable update
# ---------------------------------------------------------------------------

class TestWithVar:
    def test_returns_new_instance(self):
        env = Environment.extended({"A": "1"})
        new = env.with_var("B", "2")
        assert new is not env

    def test_new_key_added(self):
        env = Environment.extended({"A": "1"}).with_var("B", "2")
        assert env.vars == {"A": "1", "B": "2"}

    def test_existing_key_overridden(self):
        env = Environment.extended({"A": "old"}).with_var("A", "new")
        assert env.vars["A"] == "new"

    def test_original_unchanged(self):
        env = Environment.extended({"A": "1"})
        env.with_var("B", "2")
        assert "B" not in env.vars

    def test_preserves_mode_extended(self):
        env = Environment.extended().with_var("X", "1")
        assert env.is_isolated is False

    def test_preserves_mode_isolated(self):
        env = Environment.isolated({"A": "1"}).with_var("B", "2")
        assert env.is_isolated is True

    def test_chaining(self):
        env = Environment.extended().with_var("A", "1").with_var("B", "2")
        assert env.vars == {"A": "1", "B": "2"}


# ---------------------------------------------------------------------------
# resolve
# ---------------------------------------------------------------------------

class TestResolve:
    def test_isolated_returns_only_given_vars(self):
        result = Environment.isolated({"KEY": "val"}).resolve()
        assert result == {"KEY": "val"}
        assert "PATH" not in result

    def test_isolated_empty_returns_empty_dict(self):
        assert Environment.isolated({}).resolve() == {}

    def test_extended_includes_os_environ(self):
        result = Environment.extended({"SCAPHOID_UNIQUE": "yes"}).resolve()
        assert result["SCAPHOID_UNIQUE"] == "yes"
        assert "PATH" in result

    def test_extended_extras_override_os_environ(self):
        result = Environment.extended({"PATH": "overridden"}).resolve()
        assert result["PATH"] == "overridden"

    def test_extended_no_extras_equals_os_environ(self):
        assert Environment.extended().resolve() == dict(os.environ)

    def test_resolve_returns_copy(self):
        env = Environment.isolated({"A": "1"})
        d = env.resolve()
        d["NEW"] = "key"
        assert env.resolve() == {"A": "1"}


# ---------------------------------------------------------------------------
# vars property
# ---------------------------------------------------------------------------

class TestVarsProperty:
    def test_returns_explicit_vars(self):
        env = Environment.extended({"X": "1"})
        assert env.vars == {"X": "1"}

    def test_does_not_include_os_environ_for_extended(self):
        env = Environment.extended({"X": "1"})
        assert "PATH" not in env.vars

    def test_returns_copy(self):
        env = Environment.extended({"X": "1"})
        env.vars["NEW"] = "val"
        assert "NEW" not in env.vars


# ---------------------------------------------------------------------------
# immutability
# ---------------------------------------------------------------------------

class TestImmutability:
    def test_setattr_raises(self):
        env = Environment.extended()
        with pytest.raises(AttributeError):
            env._mapping = {}

    def test_setattr_arbitrary_name_raises(self):
        env = Environment.extended()
        with pytest.raises(AttributeError):
            env.new_attr = "value"


# ---------------------------------------------------------------------------
# equality and hash
# ---------------------------------------------------------------------------

class TestEquality:
    def test_equal_same_vars_and_mode(self):
        a = Environment.extended({"A": "1"})
        b = Environment.extended({"A": "1"})
        assert a == b

    def test_not_equal_different_vars(self):
        assert Environment.extended({"A": "1"}) != Environment.extended({"A": "2"})

    def test_not_equal_different_modes(self):
        assert Environment.extended({"A": "1"}) != Environment.isolated({"A": "1"})

    def test_not_equal_to_non_environment(self):
        assert Environment.extended() != {"PATH": "x"}

    def test_hashable(self):
        env = Environment.extended({"A": "1"})
        {env}  # must not raise

    def test_equal_instances_have_same_hash(self):
        a = Environment.extended({"A": "1"})
        b = Environment.extended({"A": "1"})
        assert hash(a) == hash(b)


# ---------------------------------------------------------------------------
# repr
# ---------------------------------------------------------------------------

class TestRepr:
    def test_extended_shows_mode(self):
        assert "extended" in repr(Environment.extended({"A": "1"}))

    def test_isolated_shows_mode(self):
        assert "isolated" in repr(Environment.isolated({"A": "1"}))

    def test_shows_vars(self):
        assert "hello" in repr(Environment.extended({"KEY": "hello"}))
