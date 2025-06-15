// Package goarchtest provides a fluent API for testing architectural constraints in Go applications.
//
// GoArchTest helps enforce architectural boundaries that Go's compiler cannot check.
// While Go prevents circular imports, it allows architectural violations like inner layers
// depending on outer layers. GoArchTest bridges this gap by providing executable architectural tests.
//
// # Quick Start
//
// Add architectural tests to your Go test files:
//
//	func TestArchitecture(t *testing.T) {
//	    projectPath, _ := filepath.Abs("./")
//
//	    // Test that presentation layer doesn't depend on data layer
//	    result := goarchtest.InPath(projectPath).
//	        That().
//	        ResideInNamespace("presentation").
//	        ShouldNot().
//	        HaveDependencyOn("data").
//	        GetResult()
//
//	    if !result.IsSuccessful {
//	        t.Error("Architecture violation: Presentation layer depends on data layer")
//	        for _, failingType := range result.FailingTypes {
//	            t.Logf("Violation in: %s (%s)", failingType.Name, failingType.Package)
//	        }
//	    }
//	}
//
// # Architecture Patterns
//
// GoArchTest supports predefined architecture patterns:
//
//	// Clean Architecture validation
//	cleanArch := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")
//	results := cleanArch.Validate(types)
//
// # Custom Predicates
//
// Create custom predicates for specific architectural rules:
//
//	isService := func(typeInfo *goarchtest.TypeInfo) bool {
//	    return typeInfo.IsStruct && strings.HasSuffix(typeInfo.Name, "Service")
//	}
//
//	result := types.That().
//	    WithCustomPredicate("IsService", isService).
//	    Should().
//	    ResideInNamespace("services").
//	    GetResult()
//
// # Supported Patterns
//
// - Clean Architecture
// - Hexagonal Architecture
// - Layered Architecture
// - MVC Architecture
// - DDD with Clean Architecture
// - CQRS Architecture
// - Event Sourced CQRS Architecture
//
// For more examples and documentation, visit: https://github.com/solrac97gr/goarchtest
package goarchtest

// GoArchTest is the main entry point for the architecture testing library
type GoArchTest struct {
	Types *Types
}

// New creates a new instance of GoArchTest for the specified path
func New(path string) *GoArchTest {
	return &GoArchTest{
		Types: InPath(path),
	}
}

// CheckRule executes a predefined architectural rule
// It takes a rule function that operates on the Types and returns a Result.
// This function is used to apply various architectural checks on the types defined in the Go project.
// Parameters:
//   - rule: A function that takes a pointer to Types and returns a Result
//
// Returns:
//   - *Result: Returns the result of applying the rule, which includes whether the rule passed or failed
//
// Example:
//
//	result := goarchtest.CheckRule(goarchtest.ResideInNamespace("github.com/myorg/mypackage"))
func (g *GoArchTest) CheckRule(rule func(*Types) *Result) *Result {
	return rule(g.Types)
}
