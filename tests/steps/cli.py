import json
import os
import platform
import shlex
import subprocess
import sys

from behave import given, when, then


BINARY = "./exists.exe" if sys.platform == "win32" else "./exists"
IS_WINDOWS = sys.platform == "win32"


def _run(args: list[str], extra_env: dict | None = None) -> tuple[int, str, str]:
    env = {**os.environ, **(extra_env or {})}
    proc = subprocess.run(
        [BINARY, *args],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        encoding="utf-8",
        env=env,
    )
    print(f"DEBUG stdout: {proc.stdout!r}")
    print(f"DEBUG stderr: {proc.stderr!r}")
    return proc.returncode, proc.stdout, proc.stderr

# ---------------------------------------------------------------------------
# Given
# ---------------------------------------------------------------------------

@given('the environment variable "{name}" is set to "{value}"')
def step_set_env(context, name, value):
    context.extra_env = getattr(context, "extra_env", {})
    context.extra_env[name] = value

@given('the environment variable "{name}" is empty')
def step_set_env_empty(context, name):
    context.extra_env = getattr(context, "extra_env", {})
    context.extra_env[name] = ""

# ---------------------------------------------------------------------------
# When
# ---------------------------------------------------------------------------

@when('I run exists with "{command}"')
def step_run(context, command):
    args = shlex.split(command)
    context.returncode, context.stdout, context.stderr = _run(args)


@when('I run exists with "{command}" on Windows')
def step_run_windows(context, command):
    if IS_WINDOWS:
        step_run(context, command)
    else:
        context.skip_remaining_steps = True


@when('I run exists with "{command}" on Unix')
def step_run_unix(context, command):
    if not IS_WINDOWS:
        step_run(context, command)
    elif not hasattr(context, "returncode"):
        # Windows already ran its variant; skip
        context.skip_remaining_steps = True


# ---------------------------------------------------------------------------
# Then
# ---------------------------------------------------------------------------

@then('the exit code is {code:d}')
def step_exit_code(context, code):
    if getattr(context, "skip_remaining_steps", False):
        return
    assert context.returncode == code, (
        f"Expected exit {code}, got {context.returncode}\n"
        f"stdout: {context.stdout!r}\n"
        f"stderr: {context.stderr!r}"
    )


@then('stdout contains "{text}"')
def step_stdout_contains(context, text):
    if getattr(context, "skip_remaining_steps", False):
        return
    assert text in context.stdout, (
        f"Expected stdout to contain {text!r}\n"
        f"stdout:  {context.stdout!r}"
        f"stderr:  {context.stderr!r}"
        f"context: {dir(context)!r}"
    )


@then('stderr contains "{text}"')
def step_stderr_contains(context, text):
    if getattr(context, "skip_remaining_steps", False):
        return
    assert text in context.stderr, (
        f"Expected stderr to contain {text!r}\n"
        f"stderr: {context.stderr!r}"
    )


@then('there is no output')
def step_no_output(context):
    if getattr(context, "skip_remaining_steps", False):
        return
    assert context.stdout == "" and context.stderr == "", (
        f"Expected no output\n"
        f"stdout: {context.stdout!r}\n"
        f"stderr: {context.stderr!r}"
    )


@then('the JSON output has item {index:d} with "{key}" equal to "{value}"')
def step_json_string(context, index, key, value):
    if getattr(context, "skip_remaining_steps", False):
        return
    data = json.loads(context.stdout)
    actual = data[index][key]
    assert str(actual) == value, (
        f"Expected data[{index}][{key!r}] == {value!r}, got {actual!r}"
    )


@then('the JSON output has item {index:d} with "{key}" equal to {value:w}')
def step_json_bool(context, index, key, value):
    if getattr(context, "skip_remaining_steps", False):
        return
    data = json.loads(context.stdout)
    actual = data[index][key]
    expected = {"true": True, "false": False}.get(value, value)
    assert actual == expected, (
        f"Expected data[{index}][{key!r}] == {expected!r}, got {actual!r}"
    )
