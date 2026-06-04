Feature: recipe command
  As a user
  I want to scaffold a new project from a template
  So that I can start coding without manual setup

  Scenario: Help lists available templates
    When I run alchemist "recipe --help"
    Then the exit code is 0
    And stdout contains "python"
    And stdout contains "go"
    And stdout contains "blank"

  Scenario: Rejects an unknown template immediately
    When I run alchemist "recipe unknown"
    Then the exit code is 1
    And stderr contains "unknown template"
    And stderr contains "try: python, go, blank"

  Scenario: Scaffolds a Python project from a template
    When I run alchemist "recipe python" with input
      """
      test-project
      """
    Then the exit code is 0
    And stdout contains "test-project"
    And stdout contains "ready (template: python)"

  Scenario: Scaffolds a Go project from a template
    When I run alchemist "recipe go" with input
      """
      test-project
      """
    Then the exit code is 0
    And stdout contains "test-project"
    And stdout contains "ready (template: go)"

  Scenario: Scaffolds a blank project from a template
    When I run alchemist "recipe blank" with input
      """
      test-project
      """
    Then the exit code is 0
    And stdout contains "test-project"
    And stdout contains "ready (template: blank)"
