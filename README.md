# GoArchTest

GoArchTest is a Go library for testing the architectural constraints of your Go applications. It is inspired by the .NET library [NetArchTest](https://github.com/BenMorris/NetArchTest), providing a fluent API to define and enforce architectural rules in your Go projects.

## Installation

```bash
go get github.com/solrac97gr/goarchtest
```

## Quick Start

Add architectural tests to your Go test files:

```go
func TestArchitecture(t *testing.T) {
    // Get project path
    projectPath, _ := filepath.Abs("./")
    
    // Test that presentation layer doesn't depend on data layer
    result := goarchtest.InPath(projectPath).
        That().
        ResideInNamespace("presentation").
        ShouldNot().
        HaveDependencyOn("data").
        GetResult()
        
    if !result.IsSuccessful {
        t.Error("Architecture violation: Presentation layer depends on data layer")
        for _, failingType := range result.FailingTypes {
            t.Logf("Violation in: %s (%s)", failingType.Name, failingType.Package)
        }
    }
}
```

## Why GoArchTest?

GoArchTest helps you:

- **Enforce architectural boundaries** between layers and modules
- **Detect architectural drift** early in the development process
- **Document architectural decisions** as executable tests
- **Ensure consistency** across large codebases and teams
- **Validate architectural patterns** like Clean Architecture, Hexagonal Architecture, etc.

## Detailed Usage Guide

```go
package main_test

import (
	"testing"
	"path/filepath"
	
	"github.com/solrac97gr/goarchtest"
)

func TestArchitecture(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}
	
	// Create a new Types instance for the project
	types := goarchtest.InPath(projectPath)
	
	// Test that presentation layer does not depend on data layer
	result := types.
		That().
		ResideInNamespace("presentation").
		ShouldNot().
		HaveDependencyOn("data").
		GetResult()
		
	if !result.IsSuccessful {
		t.Error("Presentation layer should not depend on data layer")
	}
}
```

### More Examples

#### Testing Layer Dependencies

```go
// Classes in the presentation should not directly reference repositories
result := types.InPath(projectPath).
    That().
    ResideInNamespace("presentation").
    ShouldNot().
    HaveDependencyOn("repository").
    GetResult()
    
if !result.IsSuccessful {
    t.Error("Presentation layer should not depend on repository layer")
}
```

#### Testing Naming Conventions

```go
// All repository implementations should have names ending with "Repository"
result := types.InPath(projectPath).
    That().
    ImplementInterface("Repository").
    Should().
    HaveNameEndingWith("Repository").
    GetResult()
    
if !result.IsSuccessful {
    t.Error("Repository implementations should end with 'Repository'")
}
```

#### Combining Predicates

```go
// Types with database dependencies should reside in the data namespace
result := types.InPath(projectPath).
    That().
    HaveDependencyOn("database/sql").
    And().
    HaveDependencyOn("gorm.io").
    Should().
    ResideInNamespace("data").
    GetResult()
    
if !result.IsSuccessful {
    t.Error("Types with database dependencies should be in data package")
}
```

### Using Predefined Architecture Patterns

GoArchTest provides helper functions for common architectural patterns:

```go
// Test clean architecture pattern
cleanArchPattern := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")
validationResults := cleanArchPattern.Validate(types)

// Report results
reporter := goarchtest.NewErrorReporter(os.Stderr)
reporter.ReportPatternValidation(validationResults)
```

### Custom Predicates

You can create custom predicates for more specific architecture rules:

```go
// Define a custom predicate function
isServiceImplementation := func(typeInfo *goarchtest.TypeInfo) bool {
    // A service implementation should:
    // 1. End with "Service"
    // 2. Be a struct
    return typeInfo.IsStruct && len(typeInfo.Name) > 7 && typeInfo.Name[len(typeInfo.Name)-7:] == "Service"
}

// Use the custom predicate
result := types.
    That().
    WithCustomPredicate("IsServiceImplementation", isServiceImplementation).
    Should().
    ResideInNamespace("application").
    GetResult()
```

### Generating Reports

GoArchTest can generate HTML or text reports of architecture test results:

```go
// Create a reporter
reporter := goarchtest.NewReporter()

// Add test results
reporter.AddResult(result1)
reporter.AddResult(result2)

// Generate and save a report
err := reporter.SaveReport("html", "architecture_report.html")
if err != nil {
    t.Fatalf("Failed to save report: %v", err)
}
```

### Visualizing Dependencies

You can generate a dependency graph in DOT format (compatible with Graphviz):

```go
// Get all types using the GetAllTypes method
allTypes := types.That().GetAllTypes()

// Create an error reporter
reporter := goarchtest.NewErrorReporter(os.Stderr)

// Save dependency graph
err = reporter.SaveDependencyGraph(allTypes, "dependency_graph.dot")
if err != nil {
    t.Logf("Failed to save dependency graph: %v", err)
}

// To generate a PNG, run:
// dot -Tpng dependency_graph.dot -o dependency_graph.png
```

For a complete working example, see the [dependency graph generation example](./examples/generate_graph) which demonstrates:

- How to analyze a project and collect type information
- How to generate a DOT file representing dependencies between packages
- How to convert the DOT file to a PNG image using Graphviz
- How to interpret the resulting dependency graph
- **How architectural violations appear in the graph** (the example includes intentional violations)

Here's an example of how a dependency graph looks for a project that violates Clean Architecture principles:

![Dependency Graph Example](./examples/generate_graph/dependency_graph.png)

In this graph:
- Each box represents a package in your project
- Arrows show dependencies between packages (arrow points to the dependency)
- Red arrows highlight architectural violations
- You can see that:
  - `application` depends on `domain` (✅ correct in Clean Architecture)
  - `infrastructure` depends on `domain` (✅ correct in Clean Architecture - outer layers can depend on inner layers)
  - `presentation` depends on `application` (✅ correct in Clean Architecture)
  - `domain` depends on `infrastructure` (❌ **violation** - inner layers should not depend on outer layers)
  - `presentation` depends on `infrastructure` (❌ **violation** - should go through application layer)

The red arrows show actual Clean Architecture violations:

1. **Domain → Infrastructure violation**: The domain layer (innermost) should never depend on the infrastructure layer (outermost). This violates the Dependency Rule which states that dependencies should only point inward.

2. **Presentation → Infrastructure violation**: The presentation layer should not directly depend on the infrastructure layer. It should only depend on the application layer, which then coordinates with infrastructure through interfaces defined in the domain.

These violations break the core principles of Clean Architecture:
- **Dependency Rule**: Dependencies must point inward toward higher-level policies
- **Separation of Concerns**: Each layer should have a single responsibility
- **Dependency Inversion**: High-level modules should not depend on low-level modules

The dependency graph visualizes package dependencies, making it easier to:
- Identify architectural violations at a glance
- Understand the overall structure of your project
- Detect circular dependencies
- Verify compliance with architectural patterns like Clean Architecture

## Implementation Details

GoArchTest uses modern Go tooling for reliable code analysis:

- Uses `golang.org/x/tools/go/packages` for robust package loading and type information
- Employs Go's type checker instead of the deprecated `ast.Package` for more accurate analysis
- Analyzes type relationships and dependencies using Go's type system
- Extracts imports, struct and interface information from Go source code

## Available Predicates

### Selectors

- `That()` - Starts a selection chain
- `And()` - Combines predicates (logical AND)
- `Or()` - Combines predicates (logical OR)
- `Should()` - Specifies that selected types should satisfy the following predicate
- `ShouldNot()` - Specifies that selected types should not satisfy the following predicate

### Filters

- `ResideInNamespace(namespace string)` - Types that reside in the specified namespace/package
- `HaveDependencyOn(dependency string)` - Types that have a dependency on the specified package
- `ImplementInterface(interfaceName string)` - Types that implement the specified interface
- `BeStruct()` - Types that are structs
- `NameMatch(pattern string)` - Types with names that match the specified regex pattern
- `HaveNameEndingWith(suffix string)` - Types with names that end with the specified suffix
- `HaveNameStartingWith(prefix string)` - Types with names that start with the specified prefix
- `ResideInDirectory(directory string)` - Types that reside in the specified directory
- `DoNotResideInNamespace(namespace string)` - Types that do not reside in the specified namespace
- `DoNotHaveDependencyOn(dependency string)` - Types that do not have a dependency on the specified package
- `WithCustomPredicate(name string, predicate CustomPredicate)` - Apply a custom predicate function

### Results

- `GetResult()` - Evaluates the predicates and returns a Result object with:
  - `IsSuccessful` - Whether the architectural test passed
  - `FailingTypes` - List of types that did not meet the criteria
  - `GetFailureDetails()` - Returns a detailed string with information about failing types

## Predefined Architecture Patterns

GoArchTest includes support for common architectural patterns:

- **Clean Architecture** - Enforces rules for domain, application, infrastructure, and presentation layers
- **Hexagonal Architecture** - Enforces rules for domain, ports, and adapters
- **Layered Architecture** - Enforces rules for a traditional n-tier architecture
- **MVC Architecture** - Enforces rules for model, view, and controller components

## Reporting and Visualization

GoArchTest includes tools for reporting and visualizing architecture test results:

- **Reporter** - Generates HTML or text reports of test results
- **ErrorReporter** - Reports errors to a specified writer (e.g., stderr)
- **Dependency Graph Generation** - Creates DOT format graphs for visualization with Graphviz

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Further Documentation

- [Real-World Examples and Best Practices](./docs/REAL_WORLD_EXAMPLES.md) - More detailed examples and best practices for using GoArchTest in real-world scenarios.
- [Frequently Asked Questions (FAQ)](./docs/FAQ.md) - Answers to common questions about using GoArchTest.
- [Error Handling Examples](./examples/error_handling/error_handling.go) - Examples of advanced error handling techniques.
- [Dependency Graph Generation Example](./examples/generate_graph) - Example showing how to generate and visualize dependency graphs for your projects.
