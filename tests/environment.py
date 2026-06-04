"""
Behave environment hooks.
Resets per-scenario state before each scenario.
"""

def before_scenario(context, scenario):
    context.returncode = None
    context.stdout = ""
    context.stderr = ""
    context.extra_env = {}
    context.skip_remaining_steps = False
