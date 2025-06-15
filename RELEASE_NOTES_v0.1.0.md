# GoArchTest v0.1.0 Release Notes

## 🎉 First Stable Release

This is the first stable release of GoArchTest, a powerful Go library for testing architectural constraints in Go applications.

## 🔧 Critical Bug Fixes

### Fixed TypeSet Mutation Bug
- **Issue**: Predicate methods were permanently modifying shared TypeSet instances, causing subsequent calls to operate on filtered results instead of original data
- **Fix**: All predicate methods now create new TypeSet instances, preserving immutability
- **Impact**: Reliable predicate chaining and consistent test results

### Fixed Namespace Matching
- **Issue**: `ResideInNamespace` failed to match relative paths like `internal/user/domain` against full module paths
- **Fix**: Enhanced matching logic with suffix and contains patterns for relative path support
- **Impact**: Proper namespace filtering for complex project structures

### Fixed Dependency Detection  
- **Issue**: `HaveDependencyOn` only performed exact/prefix matching, missing partial namespace patterns
- **Fix**: Added suffix and contains matching for comprehensive dependency detection
- **Impact**: Accurate detection of architectural violations that were previously missed

### Fixed DDD Pattern Shared Kernel Rules
- **Issue**: DDDWithCleanArchitecture pattern had overly restrictive shared kernel access rules
- **Fix**: Aligned with DDD principles - shared kernel now accessible across all layers
- **Impact**: Proper DDD pattern implementation following Domain-Driven Design principles

## ✨ Features

### Architecture Patterns
- **Clean Architecture**: Domain → Application → Infrastructure dependency flow
- **Hexagonal Architecture**: Ports and adapters pattern validation  
- **MVC Architecture**: Model-View-Controller separation
- **Layered Architecture**: Traditional n-tier architecture
- **DDD with Clean Architecture**: Bounded contexts with Clean Architecture layers
- **CQRS Architecture**: Command Query Responsibility Segregation
- **Event Sourced CQRS**: CQRS with Event Sourcing patterns

### Predicate System
- `ResideInNamespace(namespace)` - Filter by package namespace
- `HaveDependencyOn(dependency)` - Filter by import dependencies  
- `BeStruct()` / `AreInterfaces()` - Filter by type characteristics
- `HaveNameEndingWith(suffix)` / `HaveNameStartingWith(prefix)` - Filter by naming patterns
- `ImplementInterface(interface)` - Filter by interface implementation
- `Should()` / `ShouldNot()` - Positive and negative assertions
- `And()` / `Or()` - Logical operators for complex rules

### Fluent API
```go
result := goarchtest.InPath("./").
    That().
    ResideInNamespace("domain").
    ShouldNot().
    HaveDependencyOn("infrastructure").
    GetResult()
```

### Custom Predicates
```go
isRepository := func(typeInfo *goarchtest.TypeInfo) bool {
    return typeInfo.IsStruct && strings.HasSuffix(typeInfo.Name, "Repository")
}

result := types.That().
    WithCustomPredicate("IsRepository", isRepository).
    Should().
    ResideInNamespace("data").
    GetResult()
```

## 📚 Documentation

- Comprehensive package documentation with examples
- Real-world usage patterns in test examples
- Architecture pattern explanations and use cases
- Best practices for CI/CD integration

## 🧪 Test Coverage

- Clean Architecture validation example
- DDD with Clean Architecture comprehensive test suite
- Custom architecture pattern examples
- Bounded context isolation testing
- Domain-driven design pattern validation

## 🚀 Getting Started

```bash
go get github.com/solrac97gr/goarchtest@v0.1.0
```

```go
package main_test

import (
    "testing"
    "path/filepath"
    "github.com/solrac97gr/goarchtest"
)

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
        for _, failingType := range result.FailingTypes {
            t.Logf("Violation: %s in %s", failingType.Name, failingType.Package)
        }
    }
}
```

## 🔄 Breaking Changes from Alpha

- None - this release maintains backward compatibility with v0.1.0-alpha.1

## 🎯 Use Cases

- **Clean Architecture Enforcement**: Ensure proper dependency flow
- **DDD Pattern Validation**: Maintain bounded context isolation  
- **Microservice Architecture**: Validate service boundaries
- **Legacy Code Refactoring**: Prevent architectural regression
- **Team Onboarding**: Document and enforce architectural decisions
- **CI/CD Integration**: Automated architectural compliance checks

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to GoArchTest.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

**Full Changelog**: https://github.com/solrac97gr/goarchtest/compare/v0.1.0-alpha.1...v0.1.0
