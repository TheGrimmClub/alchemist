Feature: File checking
  As a user
  I want to check whether files exist
  So that I can verify required files are in place

  Scenario: Existing file
    When I run exists with "file ./examples/config.json"
    Then the exit code is 0
    And stdout contains "✓"

  Scenario: Missing file
    When I run exists with "file ./does-not-exist.json"
    Then the exit code is 1
    And stderr contains "✗"

  Scenario: Path is a directory not a file
    When I run exists with "file ./examples"
    Then the exit code is 1
    And stderr contains "not a file"

  Scenario: Multiple existing files
    When I run exists with "file ./examples/config.json ./examples/config.yaml ./examples/schema.json"
    Then the exit code is 0
    And stdout contains "3 checked"

  Scenario: JSON file passes schema validation
    When I run exists with "file ./examples/config.json --validate ./examples/schema.json"
    Then the exit code is 0
    And stdout contains "valid"

  Scenario: YAML file passes schema validation
    When I run exists with "file ./examples/config.yaml --validate ./examples/schema.json"
    Then the exit code is 0
    And stdout contains "valid"

  Scenario: Quiet mode produces no output
    When I run exists with "file ./examples/config.json --quiet"
    Then the exit code is 0
    And there is no output

  Scenario: JSON output contains expected fields
    When I run exists with "file ./examples/config.json --validate ./examples/schema.json --output-json"
    Then the exit code is 0
    And the JSON output has item 0 with "kind" equal to "file"
    And the JSON output has item 0 with "exists" equal to true
    And the JSON output has item 0 with "valid" equal to true
