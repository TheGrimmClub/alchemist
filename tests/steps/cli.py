import json
import os
import shlex

from behave import given, then, when
import subprocess

from scaphoid import EnvironmentConfig, RunConfig, BinaryExecution

BINARY = "alchemist"


def _run_cmd(context, args, stdin=None):
    """Run alchemist with the given args list, storing results on context."""
    env = EnvironmentConfig.merge(context.extra_env) if context.extra_env else None
    cfg = RunConfig(path=context.tmpdir, environment=env)
    result = BinaryExecution(cfg).run([BINARY, *args], stdin=stdin)
    context.returncode = result.exit_code
    context.stdout = result.stdout.content
    context.stderr = result.stderr.content


# ---------------------------------------------------------------------------
# Given
# ---------------------------------------------------------------------------


@given("I am in a git repository")
def step_in_git_repo(context):
    subprocess.run(["git", "init", "-q"], cwd=context.tmpdir, check=True)
    subprocess.run(["git", "config", "user.email", "test@test.com"], cwd=context.tmpdir, check=True)
    subprocess.run(["git", "config", "user.name", "Test User"], cwd=context.tmpdir, check=True)
    subprocess.run(["git", "config", "commit.gpgsign", "false"], cwd=context.tmpdir, check=True)
    with open(os.path.join(context.tmpdir, ".gitignore"), "w") as f:
        f.write(".alchemist/\n")
    subprocess.run(["git", "add", ".gitignore"], cwd=context.tmpdir, check=True)
    subprocess.run(["git", "commit", "-m", "init"], cwd=context.tmpdir, check=True)


@given('a task "{name}" is in progress')
def step_task_in_progress(context, name):
    state_dir = os.path.join(context.tmpdir, ".alchemist")
    os.makedirs(state_dir, exist_ok=True)
    state = {"task_name": name, "description": "", "created_at": "2026-01-01T00:00:00Z"}
    with open(os.path.join(state_dir, "current.json"), "w") as f:
        json.dump(state, f)


@given('a tracked file "{name}" with uncommitted changes exists')
def step_tracked_file_with_changes(context, name):
    filepath = os.path.join(context.tmpdir, name)
    with open(filepath, "w") as f:
        f.write("original content")
    subprocess.run(["git", "add", name], cwd=context.tmpdir, check=True)
    subprocess.run(["git", "commit", "-m", f"add {name}"], cwd=context.tmpdir, check=True)
    with open(filepath, "w") as f:
        f.write("modified content")


@given('an untracked file "{name}" exists')
def step_untracked_file(context, name):
    filepath = os.path.join(context.tmpdir, name)
    with open(filepath, "w") as f:
        f.write("untracked content")


@given('the environment variable "{name}" is set to "{value}"')
def step_set_env(context, name, value):
    context.extra_env[name] = value


# ---------------------------------------------------------------------------
# When
# ---------------------------------------------------------------------------


@when('I run alchemist "{command}"')
def step_run_alchemist(context, command):
    _run_cmd(context, shlex.split(command) if command else [])


@when("I run alchemist with no arguments")
def step_run_alchemist_noargs(context):
    _run_cmd(context, [])


@when('I run alchemist "{command}" with input')
def step_run_alchemist_with_input(context, command):
    stdin = context.text + "\n"
    _run_cmd(context, shlex.split(command) if command else [], stdin=stdin)


# ---------------------------------------------------------------------------
# Then
# ---------------------------------------------------------------------------


@then("the exit code is {code:d}")
def step_exit_code(context, code):
    assert context.returncode == code, (
        f"Expected exit {code}, got {context.returncode}\n"
        f"stdout: {context.stdout!r}\n"
        f"stderr: {context.stderr!r}"
    )


@then('stdout contains "{text}"')
def step_stdout_contains(context, text):
    assert text in context.stdout, (
        f"Expected stdout to contain {text!r}\n"
        f"stdout: {context.stdout!r}\n"
        f"stderr: {context.stderr!r}"
    )


@then('stderr contains "{text}"')
def step_stderr_contains(context, text):
    assert text in context.stderr, (
        f"Expected stderr to contain {text!r}\n"
        f"stderr: {context.stderr!r}"
    )


@then('the task is saved as "{name}"')
def step_task_is_saved_as(context, name):
    state_file = os.path.join(context.tmpdir, ".alchemist", "current.json")
    assert os.path.exists(state_file), "State file .alchemist/current.json does not exist"
    with open(state_file) as f:
        state = json.load(f)
    assert state["task_name"] == name, (
        f"Expected task_name {name!r}, got {state['task_name']!r}"
    )


@then("no task is in progress")
def step_no_task(context):
    state_file = os.path.join(context.tmpdir, ".alchemist", "current.json")
    assert not os.path.exists(state_file), (
        "Expected no active task, but .alchemist/current.json exists"
    )


@then('the file "{name}" does not exist')
def step_file_not_exist(context, name):
    filepath = os.path.join(context.tmpdir, name)
    assert not os.path.exists(filepath), f"Expected {name!r} to be deleted but it still exists"
