Feature: stash command
  As a user
  I want to pause my work and set it aside
  So that I can switch context without losing uncommitted progress

  Scenario: Help shows the --resume flag
    When I run alchemist "stash --help"
    Then the exit code is 0
    And stdout contains "--resume"

  Scenario: Fails outside a git repository
    When I run alchemist "stash"
    Then the exit code is 1
    And stderr contains "not a git repository"

  Scenario: Stash resume fails outside a git repository
    When I run alchemist "stash --resume"
    Then the exit code is 1
    And stderr contains "not a git repository"

  Scenario: Shows confirmation prompt when run in a repository
    Given I am in a git repository
    When I run alchemist "stash" with input
      """
      n
      """
    Then the exit code is 0
    And stdout contains "Stash: pause your work"

  Scenario: Cancels stash when the user declines
    Given I am in a git repository
    When I run alchemist "stash" with input
      """
      n
      """
    Then the exit code is 0
    And stdout contains "Cancelled"
