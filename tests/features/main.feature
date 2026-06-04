Feature: Main command
  As a user
  I want to learn more about the `alchemist` command and its capabilities
  So that I can use it effectively in my projects

  Scenario: Help command shows usage information
    When I run exists with "--help"
    Then the exit code is 0
    And stdout contains "Usage:"
    And stdout contains "Available Commands:"
