# Contributing to GoArchTest

Thank you for your interest in contributing to GoArchTest! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to maintain a respectful and inclusive community.

## How to Contribute

There are many ways to contribute to GoArchTest:

- Reporting bugs
- Suggesting enhancements
- Improving documentation
- Adding new features
- Fixing bugs
- Writing tests

## Development Process

1. Fork the repository
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes
5. Run tests to ensure everything works
6. Submit a pull request

### Setting Up Development Environment

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/goarchtest.git
cd goarchtest

# Add the original repository as a remote
git remote add upstream https://github.com/solrac97gr/goarchtest.git

# Create a new branch for your changes
git checkout -b feature/your-feature-name
```

### Running Tests

```bash
go test ./...
```

## Pull Request Process

1. Update the README.md and documentation with details of changes if appropriate
2. Update examples if you've added or changed functionality
3. Make sure all tests pass
4. Your pull request will be reviewed by the maintainers
5. Once approved, your changes will be merged

## Coding Standards

- Follow Go's official [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use [gofmt](https://golang.org/cmd/gofmt/) to format your code
- Add comments explaining non-obvious code
- Write tests for new functionality
- Follow the existing code style and structure

## Architecture

GoArchTest is organized into several components:

- **Core Types and Predicates**: The foundation for architecture testing
- **Pattern Validators**: Pre-defined architectural pattern rules
- **Reporting Tools**: Error reporting and visualization

When adding new features, consider which component your changes fit into.

## Feature Requests

Feature requests are welcome! To suggest a new feature:

1. Check the issue tracker to see if it's already been suggested
2. If not, create a new issue with the "enhancement" label
3. Clearly describe the feature and why it would be valuable

## Bug Reports

When reporting bugs, please include:

1. Steps to reproduce the bug
2. Expected behavior
3. Actual behavior
4. Go version and environment information
5. Code examples if applicable

## License

By contributing to GoArchTest, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).

## Questions?

If you have questions about contributing, feel free to open an issue or contact the maintainers.
