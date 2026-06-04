Feature: start command
  As a user
  I want to name the task before I write any code
  So that my commit message is decided upfront

  Scenario: Help shows the command's purpose
    When I run alchemist "start --help"
    Then the exit code is 0
    And stdout contains "Records the task name"

  Scenario: Fails outside a git repository
    When I run alchemist "start"
    Then the exit code is 1
    And stderr contains "not a git repository"

  Scenario: Saves a new task and confirms it
    Given I am in a git repository
    When I run alchemist "start" with input
      """
      My Task
      My description
      """
    Then the exit code is 0
    And stdout contains "Task started: My Task"
    And the task is saved as "My Task"

  Scenario: Detects an already-started task and offers to replace it
    Given I am in a git repository
    And a task "Old Task" is in progress
    When I run alchemist "start" with input
      """
      n
      """
    Then the exit code is 0
    And stdout contains "A task is already in progress"
    And stdout contains "Old Task"

  Scenario: Replaces an existing task when the user confirms
    Given I am in a git repository
    And a task "Old Task" is in progress
    When I run alchemist "start" with input
      """
      y
      New Task

      """
    Then the exit code is 0
    And stdout contains "Task started: New Task"
    And the task is saved as "New Task"
