"""
Behave environment hooks.
Sets up and tears down per-scenario state and a temporary working directory.
"""

import os
import shutil
import tempfile


def before_scenario(context, scenario):
    context.returncode = None
    context.stdout = ""
    context.stderr = ""
    context.extra_env = {}
    context.tmpdir = tempfile.mkdtemp(prefix="alchemist_test_")


def after_scenario(context, scenario):
    if hasattr(context, "tmpdir") and os.path.isdir(context.tmpdir):
        shutil.rmtree(context.tmpdir)
