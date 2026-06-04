Feature: brew command
  As a user
  I want to commit my work with a meaningful message
  So that my git history reflects intentional decisions

  Scenario: Help shows the command's purpose
    When I run alchemist "brew --help"
    Then the exit code is 0
    And stdout contains "Stages your changes"

  Scenario: Fails outside a git repository
    When I run alchemist "brew"
    Then the exit code is 1
    And stderr contains "not a git repository"

  Scenario: Fails when no task has been started
    Given I am in a git repository
    When I run alchemist "brew"
    Then the exit code is 1
    And stderr contains "no task in progress"

  Scenario: Reports nothing to commit when the working tree is clean
    Given I am in a git repository
    And a task "My Feature" is in progress
    When I run alchemist "brew"
    Then the exit code is 0
    And stdout contains "Nothing to commit"

  Scenario: Reads the task state saved by start
    Given I am in a git repository
    And a task "State Check Task" is in progress
    When I run alchemist "brew"
    Then the exit code is 0
    And stdout contains "Nothing to commit"
