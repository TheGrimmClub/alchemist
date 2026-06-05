Feature: bottle command
  As a user
  I want to finalize my work with an optional tag and push
  So that my commits are published and the task state is cleared

  Scenario: Help shows the command's purpose
    When I run alchemist "bottle --help"
    Then the exit code is 0
    And stdout contains "Pushes your committed work"

  Scenario: Fails outside a git repository
    When I run alchemist "bottle"
    Then the exit code is 1
    And stderr contains "not a git repository"

  Scenario: Shows push confirmation prompt
    Given I am in a git repository
    When I run alchemist "bottle" with input
      """
      n
      """
    Then the exit code is 0
    And stdout contains "Bottle: push your work"

  Scenario: Cancels push when the user declines
    Given I am in a git repository
    When I run alchemist "bottle" with input
      """
      n
      """
    Then the exit code is 0
    And stdout contains "Cancelled"
