Feature: Environment variable checking
  As a user
  I want to check whether environment variables are set
  So that I can verify my runtime configuration

  Scenario: Set environment variable
    When I run exists with "environment PATH"
    Then the exit code is 0
    And stdout contains "✓"

  Scenario: Missing environment variable
    When I run exists with "environment THIS_VAR_DOES_NOT_EXIST_XYZ"
    Then the exit code is 1
    And stderr contains "✗"

  Scenario: Multiple set environment variables
    When I run exists with "environment PATH TEMP" on Windows
    When I run exists with "environment PATH HOME" on Unix
    Then the exit code is 0
    And stdout contains "2 checked"

  Scenario: Non-empty check passes for set variable
    When I run exists with "environment PATH --non-empty"
    Then the exit code is 0

  Scenario: Non-empty check fails for empty variable
    Given the environment variable "EXISTS_TEST_EMPTY" is empty
    When I run exists with "environment EXISTS_TEST_EMPTY --non-empty"
    Then the exit code is 1
    And stderr contains "empty"

  Scenario: Quiet mode produces no output
    When I run exists with "environment PATH --quiet"
    Then the exit code is 0
    And there is no output

  Scenario: JSON output contains expected fields
    When I run exists with "environment PATH --output-json"
    Then the exit code is 0
    And the JSON output has item 0 with "kind" equal to "env"
    And the JSON output has item 0 with "exists" equal to true
