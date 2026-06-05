Feature: Root command
  As a user
  I want to see what alchemist can do
  So that I can choose the right subcommand

  Scenario: Help shows usage and all subcommands
    When I run alchemist "--help"
    Then the exit code is 0
    And stdout contains "Usage:"
    And stdout contains "Available Commands:"
    And stdout contains "start"
    And stdout contains "brew"
    And stdout contains "bottle"
    And stdout contains "discard"
    And stdout contains "stash"
    And stdout contains "clean"
    And stdout contains "recipe"

  Scenario: Running without arguments shows help
    When I run alchemist with no arguments
    Then the exit code is 0
    And stdout contains "Usage:"
