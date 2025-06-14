# 🚀 GoArchTest v0.1.0-alpha.1 - First Pre-release

Welcome to the first pre-release of **GoArchTest**, a Go library for testing architectural constraints in your Go applications! This is an alpha version intended for early adopters and feedback collection.

## 🎯 What is GoArchTest?

GoArchTest helps you enforce architectural boundaries that Go's compiler cannot check. While Go prevents circular imports, it allows architectural violations like inner layers depending on outer layers. GoArchTest bridges this gap by providing executable architectural tests.

## ✨ Key Features

### 🏗️ Architecture Patterns
- **Clean Architecture** - Domain, Application, Infrastructure, Presentation layers
- **Hexagonal Architecture** - Domain, Ports, Adapters
- **Layered Architecture** - Traditional n-tier architectures
- **MVC Architecture** - Model, View, Controller
- **DDD + Clean Architecture** - Domain-Driven Design with Clean Architecture
- **CQRS Architecture** - Command Query Responsibility Segregation
- **Event Sourced CQRS** - Event Sourcing with CQRS

### 🔍 Powerful Predicates
- Namespace/package filtering
- Dependency analysis
- Interface implementation checks
- Name pattern matching (regex, prefix, suffix)
- Struct/Interface type filtering
- Custom predicate support
- Directory-based filtering

### 📊 Reporting & Visualization
- HTML and text report generation
- DOT format dependency graphs
- Graphviz integration for visualization
- Detailed error reporting

### 🔧 CI/CD Ready
- GitHub Actions integration
- Automated architecture testing
- Multi-environment support

## 🚀 Quick Start

```bash
# Install the pre-release
go get github.com/solrac97gr/goarchtest@v0.1.0-alpha.1
```

```go
func TestArchitecture(t *testing.T) {
    projectPath, _ := filepath.Abs("./")
    
    // Test Clean Architecture rules
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

## 📚 Examples Included

- **Clean Architecture Validation** - Complete example with layers
- **Custom Predicates** - Extensible rule definitions
- **Naming Conventions** - Service, handler, repository patterns
- **Dependency Graphs** - Visual architecture analysis

## 🎯 Perfect For

- **Teams adopting Clean Architecture** - Enforce proper layer boundaries
- **Microservice architectures** - Validate service boundaries
- **Large codebases** - Prevent architectural drift
- **Code reviews** - Automated architecture compliance
- **Onboarding** - Document architectural decisions as tests

## ⚠️ Pre-release Notice

This is an **alpha pre-release**:
- API may change based on feedback
- Intended for early adopters and testing
- Please report issues and share feedback
- Documentation and examples are evolving

## 🔗 Resources

- **Documentation**: [README.md](https://github.com/solrac97gr/goarchtest/blob/main/README.md)
- **Examples**: [examples/](https://github.com/solrac97gr/goarchtest/tree/main/examples)
- **FAQ**: [docs/FAQ.md](https://github.com/solrac97gr/goarchtest/blob/main/docs/FAQ.md)
- **Contributing**: [CONTRIBUTING.md](https://github.com/solrac97gr/goarchtest/blob/main/CONTRIBUTING.md)

## 🤝 Get Involved

We're actively seeking feedback on:
- API design and usability
- Additional architecture patterns
- Performance with large codebases
- Integration scenarios
- Documentation clarity

**Found a bug?** Open an issue
**Have a feature idea?** Start a discussion
**Want to contribute?** Check our contributing guide

## 🙏 Acknowledgments

Inspired by [NetArchTest](https://github.com/BenMorris/NetArchTest) for .NET, bringing similar architectural testing capabilities to the Go ecosystem.

---

**Try it today and help shape the future of architectural testing in Go!** 🎉
