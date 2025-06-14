# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0-alpha.1] - 2025-06-14

### 🚀 Added

#### Core Features
- **Fluent API** for defining and enforcing architectural rules
- **Package Analysis** using Go's type checker and AST parsing
- **Architecture Testing** with comprehensive predicate system
- **Dependency Analysis** for detecting architectural violations

#### Architecture Patterns
- **Clean Architecture** validation with domain, application, infrastructure, and presentation layers
- **Hexagonal Architecture** validation with domain, ports, and adapters
- **Layered Architecture** validation for traditional n-tier architectures
- **MVC Architecture** validation for model, view, and controller components
- **DDD with Clean Architecture** validation for Domain-Driven Design
- **CQRS Architecture** validation for Command Query Responsibility Segregation
- **Event Sourced CQRS Architecture** validation for Event Sourcing with CQRS

#### Predicates & Filters
- `ResideInNamespace(namespace)` - Filter types by namespace/package
- `HaveDependencyOn(dependency)` - Filter types with specific dependencies
- `ImplementInterface(interfaceName)` - Filter types implementing interfaces
- `BeStruct()` - Filter struct types
- `AreInterfaces()` - Filter interface types
- `NameMatch(pattern)` - Filter types by regex pattern
- `HaveNameMatching(pattern)` - Alias for NameMatch
- `HaveNameEndingWith(suffix)` - Filter types by name suffix
- `HaveNameStartingWith(prefix)` - Filter types by name prefix
- `ResideInDirectory(directory)` - Filter types by directory
- `DoNotResideInNamespace(namespace)` - Exclude types from namespace
- `DoNotHaveDependencyOn(dependency)` - Exclude types with dependencies
- `WithCustomPredicate(name, predicate)` - Apply custom predicate functions

#### Logical Operators
- `And()` - Combine predicates with logical AND
- `Or()` - Combine predicates with logical OR
- `Should()` - Specify positive conditions
- `ShouldNot()` - Specify negative conditions

#### Reporting & Visualization
- **HTML Report Generation** for architecture test results
- **Text Report Generation** for command-line output
- **DOT Format Dependency Graphs** compatible with Graphviz
- **Error Reporting** with detailed violation information
- **Dependency Visualization** with PNG generation support

#### Examples & Documentation
- **Clean Architecture Example** with complete validation tests
- **Custom Predicate Examples** showing extensibility
- **Architecture Pattern Definitions** for common patterns
- **Dependency Graph Generation** examples
- **Real-world Usage Scenarios** documentation

#### CI/CD Integration
- **GitHub Actions Workflow** for automated testing
- **Multi-environment Testing** across Go versions
- **Dependency Caching** for faster builds
- **Comprehensive Test Suite** validation

### 🔧 Technical Details

#### Dependencies
- `golang.org/x/tools` v0.34.0 for Go package analysis
- `golang.org/x/mod` v0.25.0 for module handling
- `golang.org/x/sync` v0.15.0 for synchronization primitives

#### Go Version Support
- Requires Go 1.24.1 or later
- Tested on Go 1.24.x

#### Package Structure
- Core library in root package
- Examples in `examples/` directory
- Test cases in `test/` directory
- Documentation in `docs/` directory

### 📋 API Reference

#### Main Entry Points
```go
// Create Types instance for a project
types := goarchtest.InPath(projectPath)

// Start filtering chain
result := types.That().
    ResideInNamespace("domain").
    ShouldNot().
    HaveDependencyOn("infrastructure").
    GetResult()
```

#### Architecture Patterns
```go
// Use predefined patterns
cleanArch := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")
results := cleanArch.Validate(types)
```

#### Custom Predicates
```go
// Define custom rules
isService := func(t *goarchtest.TypeInfo) bool {
    return t.IsStruct && strings.HasSuffix(t.Name, "Service")
}

result := types.That().
    WithCustomPredicate("IsService", isService).
    Should().
    ResideInNamespace("services").
    GetResult()
```

### 🎯 Use Cases

#### Layer Validation
- Prevent inner layers from depending on outer layers
- Enforce proper dependency direction in Clean Architecture
- Validate hexagonal architecture port/adapter boundaries

#### Naming Conventions
- Ensure services end with "Service"
- Validate handler naming patterns
- Check interface naming conventions

#### Dependency Management
- Prevent business logic from depending on frameworks
- Ensure data access components use appropriate libraries
- Validate cross-cutting concern separation

#### Team Consistency
- Enforce architectural decisions across teams
- Document architectural rules as executable tests
- Catch architectural drift early in CI/CD

### 🔮 Future Plans

#### Upcoming Features (v0.2.0)
- Performance optimizations for large codebases
- Additional predefined architecture patterns
- Enhanced reporting formats (JSON, XML)
- Integration with popular IDEs

#### Long-term Roadmap
- Plugin system for custom architecture patterns
- Integration with architectural documentation tools
- Support for microservice architecture validation
- Advanced dependency analysis algorithms

### 📦 Installation

```bash
# Install pre-release version
go get github.com/solrac97gr/goarchtest@v0.1.0-alpha.1

# Install latest stable (when available)
go get github.com/solrac97gr/goarchtest
```

### 🤝 Contributing

This is a pre-release version. We welcome:
- Bug reports and feature requests
- Documentation improvements
- Example contributions
- Architecture pattern suggestions

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### ⚠️ Pre-release Notice

This is an alpha pre-release intended for:
- Early adopters and testing
- Feedback collection
- API refinement

The API may change in future releases based on community feedback.
