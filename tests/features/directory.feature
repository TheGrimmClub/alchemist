Feature: Directory checking
  As a user
  I want to check whether directories exist
  So that I can verify my project structure

  Scenario: Existing directory
    When I run exists with "directory ./examples"
    Then the exit code is 0
    And stdout contains "✓"

  Scenario: Missing directory
    When I run exists with "directory ./does-not-exist"
    Then the exit code is 1
    And stderr contains "✗"

  Scenario: Path is a file not a directory
    When I run exists with "directory ./examples/config.json"
    Then the exit code is 1
    And stderr contains "not a directory"

  Scenario: Multiple existing directories
    When I run exists with "directory ./examples ./cmd ./internal"
    Then the exit code is 0
    And stdout contains "3 checked"

  Scenario: Mixed existing and missing directories
    When I run exists with "directory ./examples ./does-not-exist"
    Then the exit code is 1
    And stdout contains "2 checked"

  Scenario: Quiet mode produces no output
    When I run exists with "directory ./examples --quiet"
    Then the exit code is 0
    And there is no output

  Scenario: JSON output contains expected fields
    When I run exists with "directory ./examples --output-json"
    Then the exit code is 0
    And the JSON output has item 0 with "kind" equal to "directory"
    And the JSON output has item 0 with "exists" equal to true
