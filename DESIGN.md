# Design

A TUI (Terminal User Interface) custom to my current workflows at Cinergia.

## Glossary

- Task : a unit of work with title, steps (checkable) and text area.

## Features

**Template**:

```gherkin
Scenario: <scenario name>
    Given <precondition>
        And <additional precondition>
When <action>
    Then <expected outcome>
        And <additional expected outcome>
```

**Feature List**:

- [Task Creation](#feature-task-creation)
- [Task Status Editing](#feature-task-status-editing)
- [Task Editing](#feature-task-editing)
- [Task Steps Editing](#feature-task-steps-editing)
- [Task Deletion](#feature-task-deletion)
- [Task Archiving](#feature-task-archiving)

---

### Feature: Task Creation

```gherkin
Scenario: Create new task
    Given I am on the task input
    When I write the task name
        And I press Intro
    Then a new task is created at the main panel
        And has status is "New"
        And has no steps

Scenario: Task name exists
    Given I am on the task input
    When I write the task name
        And it already exists
        And I press intro
    Then an error message appears under the input field
        And the error is red
        And the error states "name already exists"

Scenario: Task name is empty
    Given I am on the task input
    When I leave the name blank
        And press intro
    Then an error message appears under the input
        And the error is red
        And the error states "name can't be empty"
```

### Feature: Task Status Editing

```gherkin
Scenario: Edit task status 
    Given I am viewing a task card
        And it is not archived
        And it is not removed
    When I change the status to <status>
    Then the task status is uptaded to <status>

    Examples:
        | Status            | 
        | New               |
        | In Progress       |
        | Ready to Validate |
        | Done              |
```

### Feature: Task Editing

```gherkin
Scenario: Edit task name
    Given I am viewing a task card
        And the task is not archived
    When I edit the task name
        And I confirm the change
    Then the task name is updated

Scenario: Edit task description
    Given I am viewing a task card
        And the task is not archived
    When I edit the task description
        And I confirm the change
    Then the task description is updated
```

### Feature: Task Steps Editing

```gherkin
Scenario: Add step to task
    Given I am viewing a task card
    When I write the step text in the step input
        And I press Enter
    Then a new step is added to the task
        And the step is marked as incomplete

Scenario: Edit step text
    Given I am viewing a task card
        And the task has at least one step
    When I edit a step's text
        And I confirm the change
    Then the step text is updated

Scenario: Mark step as complete
    Given I am viewing a task card
        And the task has a step that is not complete
    When I mark the step as complete
    Then the step is marked as complete

Scenario: Mark step as incomplete
    Given I am viewing a task card
        And the task has a step that is complete
    When I mark the step as incomplete
    Then the step is marked as incomplete

Scenario: Delete step
    Given I am viewing a task card
        And the task has at least one step
    When I delete a step
    Then the step is removed from the task
```

---

### Feature: Task Deletion

```gherkin
Scenario: Delete task
    Given I am viewing a task card
    When I delete the task
    Then the task is removed from the task list
```

---

### Feature: Task Archiving

```gherkin
Scenario: Archive task
    Given I am viewing a task card
        And the task is not archived
    When I archive the task
    Then the task is removed from the main panel
        And it appears in the archive

Scenario: Restore archived task
    Given I am viewing an archived task
    When I restore the task
    Then the task appears in the main panel
        And it is no longer in the archive
```
