Feature: clean command
  As a user
  I want to remove untracked files from my project
  So that build artifacts and stray files don't accumulate

  Scenario: Help shows the command's purpose
    When I run alchemist "clean --help"
    Then the exit code is 0
    And stdout contains "Deletes files"

  Scenario: Fails outside a git repository
    When I run alchemist "clean"
    Then the exit code is 1
    And stderr contains "not a git repository"

  Scenario: Reports nothing to clean in a repository with no untracked files
    Given I am in a git repository
    When I run alchemist "clean"
    Then the exit code is 0
    And stdout contains "Nothing to clean"
