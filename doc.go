/*
Package goarchtest provides a fluent API for testing architectural constraints in Go applications.

# Overview

GoArchTest helps enforce architectural boundaries that Go's compiler cannot check.
While Go prevents circular imports, it allows architectural violations like inner layers
depending on outer layers. GoArchTest bridges this gap by providing executable architectural tests.

# Key Features

  - Fluent API for readable architectural tests
  - Predefined architecture patterns (Clean, Hexagonal, MVC, etc.)
  - Extensive predicate system for type filtering
  - Dependency analysis and violation detection
  - Visual dependency graph generation
  - CI/CD integration support

# Quick Start

Add architectural tests to your Go test files:

	func TestArchitecture(t *testing.T) {
	    projectPath, _ := filepath.Abs("./")

	    result := goarchtest.InPath(projectPath).
	        That().
	        ResideInNamespace("domain").
	        ShouldNot().
	        HaveDependencyOn("infrastructure").
	        GetResult()

	    if !result.IsSuccessful {
	        t.Error("Domain layer should not depend on infrastructure")
	    }
	}

# Architecture Patterns

GoArchTest supports several predefined architecture patterns:

	// Clean Architecture
	pattern := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")
	results := pattern.Validate(types)

	// Hexagonal Architecture
	pattern := goarchtest.HexagonalArchitecture("domain", "ports", "adapters")
	results := pattern.Validate(types)

	// MVC Architecture
	pattern := goarchtest.MVCArchitecture("models", "views", "controllers")
	results := pattern.Validate(types)

# Available Predicates

The predicate system allows flexible filtering and testing of types:

## Type Filters

  - ResideInNamespace(namespace) - Filter by package namespace
  - BeStruct() - Filter struct types only
  - AreInterfaces() - Filter interface types only
  - HaveNameEndingWith(suffix) - Filter by type name suffix
  - HaveNameStartingWith(prefix) - Filter by type name prefix
  - HaveNameMatching(pattern) - Filter by regex pattern

## Dependency Analysis

  - HaveDependencyOn(dependency) - Filter types with specific dependencies
  - DoNotHaveDependencyOn(dependency) - Filter types without dependencies
  - ImplementInterface(interfaceName) - Filter types implementing interfaces

## Logical Operators

  - And() - Combine predicates with logical AND
  - Or() - Combine predicates with logical OR
  - Should() - Specify positive conditions
  - ShouldNot() - Specify negative conditions (negation)

# Custom Predicates

Create custom rules for specific architectural constraints:

	isRepository := func(typeInfo *goarchtest.TypeInfo) bool {
	    return typeInfo.IsStruct && strings.HasSuffix(typeInfo.Name, "Repository")
	}

	result := types.That().
	    WithCustomPredicate("IsRepository", isRepository).
	    Should().
	    ResideInNamespace("data").
	    GetResult()

# Reporting and Visualization

Generate reports and visualizations of your architecture:

	// Generate dependency graph
	allTypes := types.That().GetAllTypes()
	reporter := goarchtest.NewErrorReporter(os.Stderr)
	reporter.SaveDependencyGraph(allTypes, "dependencies.dot")

	// Generate HTML report
	reporter := goarchtest.NewReporter()
	reporter.AddResult(result)
	reporter.SaveReport("html", "architecture_report.html")

# Real-World Examples

## Clean Architecture Validation

	func TestCleanArchitecture(t *testing.T) {
	    types := goarchtest.InPath("./")

	    // Domain should not depend on infrastructure
	    result := types.That().
	        ResideInNamespace("domain").
	        ShouldNot().
	        HaveDependencyOn("infrastructure").
	        GetResult()

	    assert.True(t, result.IsSuccessful)
	}

## Service Layer Validation

	func TestServiceLayer(t *testing.T) {
	    types := goarchtest.InPath("./")

	    // All services should end with "Service"
	    result := types.That().
	        ResideInNamespace("services").
	        Should().
	        HaveNameEndingWith("Service").
	        GetResult()

	    assert.True(t, result.IsSuccessful)
	}

## Naming Convention Validation

	func TestNamingConventions(t *testing.T) {
	    types := goarchtest.InPath("./")

	    // Repositories should be in data package
	    result := types.That().
	        HaveNameEndingWith("Repository").
	        Should().
	        ResideInNamespace("data").
	        GetResult()

	    assert.True(t, result.IsSuccessful)
	}

# CI/CD Integration

Include architectural tests in your continuous integration:

	name: Architecture Tests
	on: [push, pull_request]
	jobs:
	  test:
	    runs-on: ubuntu-latest
	    steps:
	      - uses: actions/checkout@v4
	      - uses: actions/setup-go@v4
	        with:
	          go-version: '1.24.x'
	      - run: go test -v ./...

# Use Cases

GoArchTest is perfect for:

  - Teams adopting Clean Architecture or Hexagonal Architecture
  - Large codebases where architectural drift is a concern
  - Microservice architectures requiring boundary validation
  - Code reviews requiring architectural compliance
  - Onboarding new developers with architectural documentation
  - Enforcing coding standards and conventions

# Why GoArchTest?

While Go's compiler prevents circular imports (A→B→A), it allows architectural violations:

	// Go compiler: ✅ COMPILES FINE
	// Clean Architecture: ❌ VIOLATION
	package domain
	import "infrastructure" // Inner layer importing outer layer

	// GoArchTest catches this:
	result := types.That().
	    ResideInNamespace("domain").
	    ShouldNot().
	    HaveDependencyOn("infrastructure"). // ❌ Test fails!
	    GetResult()

GoArchTest bridges the gap between what Go's compiler enforces and what good
software architecture requires.

For more examples and detailed documentation, visit:
https://github.com/solrac97gr/goarchtest
*/
package goarchtest
