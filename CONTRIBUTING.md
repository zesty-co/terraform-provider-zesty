# Contributing to Terraform Provider for Zesty

- [Contributing to Terraform Provider for Zesty](#contributing-to-terraform-provider-for-zesty)
    - [Pull Requests](#pull-requests)
        - [Titles](#titles)
        - [Prefixes](#prefixes)
        - [Functional Changes](#functional-changes)
        - [Non-functional changes](#non-functional-changes)
    - [Coding Standards](#coding-standards)
    - [Testing](#testing)

## Pull Requests

#### Titles

Your PRs titles should have the following format:

```
<type>: <description> (Issue-Reference)
```

For example:

```
fix: reduced error rate and increased stability (#1337)
```

##### Prefixes

Note that we are using semver, thus functional changes will increase
the minor version, whie non-functional changes will increase the patch
version. Breaking changes will increase the major version.

###### Functional Changes

- `feat` - A new feature or a core change to an existing one.

###### Non-functional changes

- `fix` - A fix to a bug or a change in implementation that does not change behavior.
- `test` - An addition or change to the testing suite.
- `refactor` - A change that does not impact behavior but is done for engineering reasons.
- `docs` - An update or addition to the documentation.
- `ci` - A change to the GitHub actions.

## Coding Standards

The builtin Go coding standards apply. There are several commands in the project's
Makefile that helps achieve that.

- Before running any tool, make sure you have the installed tools, this can be done easily using `make tools`.

- Running `make format` will make sure your code is formatted correctly.

- Running `make lint` will make sure all the code linters are run against your code changes.
    - Also, `make lint-fix` will try to fix any issues that have auto-fix capability in the linter.

- Running `make generate` will generate all code generated content according to your changes.

- Running `make docs` will generate the auto-generated documentation used in the tf registry.

## Testing

Make sure code changes or new features are being tested and the relevant code is
covered. Running `make test` will run the test and produce a `coverage.html` to easily
browse the code coverage.
