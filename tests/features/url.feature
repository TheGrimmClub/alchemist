Feature: URL checking
  As a user
  I want to check whether URLs are reachable
  So that I can verify external dependencies are available

  Scenario: Reachable URL
    When I run exists with "url https://example.com"
    Then the exit code is 0
    And stdout contains "✓"

  Scenario: Correct expected status code
    When I run exists with "url https://example.com --expect 200"
    Then the exit code is 0

  Scenario: Wrong expected status code
    When I run exists with "url https://example.com --expect 418"
    Then the exit code is 1
    And stderr contains "✗"

  Scenario: Multiple reachable URLs
    When I run exists with "url https://example.com https://example.org"
    Then the exit code is 0
    And stdout contains "2 checked"

  Scenario: Unreachable URL
    When I run exists with "url https://this.does.not.exist.invalid"
    Then the exit code is 1

  Scenario: Custom timeout is accepted
    When I run exists with "url https://example.com --timeout 10s"
    Then the exit code is 0

  Scenario: No redirect following
    When I run exists with "url https://example.com --follow=false"
    Then the exit code is 0

  Scenario: JSON output contains expected fields
    When I run exists with "url https://example.com --output-json"
    Then the exit code is 0
    And the JSON output has item 0 with "kind" equal to "url"
    And the JSON output has item 0 with "exists" equal to true
