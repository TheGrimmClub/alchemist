Feature: discard command
  As a user
  I want to throw away uncommitted changes to tracked files
  So that I can start over without switching branches

  Scenario: Help shows the command's purpose
    When I run alchemist "discard --help"
    Then the exit code is 0
    And stdout contains "Reverts tracked files"

  Scenario: Fails outside a git repository
    When I run alchemist "discard"
    Then the exit code is 1
    And stderr contains "not a git repository"

  Scenario: Reports nothing to discard in a clean repository
    Given I am in a git repository
    When I run alchemist "discard"
    Then the exit code is 0
    And stdout contains "Nothing to discard"

  Scenario: Shows changed files and cancels when the user declines
    Given I am in a git repository
    And a tracked file "hello.txt" with uncommitted changes exists
    When I run alchemist "discard" with input
      """
      n
      """
    Then the exit code is 0
    And stdout contains "These changes will be permanently thrown away"
    And stdout contains "hello.txt"
    And stdout contains "Cancelled"

  Scenario: Discards changes when the user confirms
    Given I am in a git repository
    And a tracked file "hello.txt" with uncommitted changes exists
    When I run alchemist "discard" with input
      """
      y
      """
    Then the exit code is 0
    And stdout contains "Changes discarded"
