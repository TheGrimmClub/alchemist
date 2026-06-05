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

  Scenario: Shows untracked files and cancels when the user declines
    Given I am in a git repository
    And an untracked file "build.out" exists
    When I run alchemist "clean" with input
      """
      n
      """
    Then the exit code is 0
    And stdout contains "build.out"
    And stdout contains "Cancelled"

  Scenario: Deletes untracked files when the user confirms
    Given I am in a git repository
    And an untracked file "build.out" exists
    When I run alchemist "clean" with input
      """
      y
      """
    Then the exit code is 0
    And stdout contains "removed"
    And the file "build.out" does not exist
