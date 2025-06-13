# Frequently Asked Questions (FAQ)

## General Questions

### What is GoArchTest?

GoArchTest is a Go library for testing architectural constraints in Go applications. It provides a fluent API to define rules about how your code should be structured and helps ensure that your codebase adheres to those rules.

### Why do I need GoArchTest when Go already prevents import cycles?

This is a common and excellent question! While Go's compiler prevents **circular dependencies**, it doesn't prevent **architectural violations**. Here's the key difference:

**Go prevents this (import cycle):**
```go
package A
import "B"  // A → B

package B  
import "A"  // B → A (ERROR: import cycle)
```

**Go ALLOWS this (but it violates Clean Architecture):**
```go
package domain          // Inner layer
import "infrastructure" // Depending on outer layer ❌

package infrastructure  // Outer layer  
import "domain"         // Depending on inner layer ✅
```

Both packages compile fine in Go (no import cycle), but the first violates architectural principles. GoArchTest catches violations like:
- Domain layer importing infrastructure
- Presentation layer skipping application to access infrastructure directly
- Business logic depending on HTTP frameworks
- Cross-domain contamination between bounded contexts

**Example of what GoArchTest catches:**
```go
// This compiles in Go but breaks Clean Architecture
package domain
import "myapp/infrastructure/database"  // ❌ Architecture violation

// GoArchTest detects this:
func TestArchitecture(t *testing.T) {
    result := goarchtest.InPath("./").
        That().ResideInNamespace("domain").
        ShouldNot().HaveDependencyOn("infrastructure").
        GetResult()  // ❌ Fails - catches the violation!
}
```

So GoArchTest provides **architectural enforcement** beyond Go's **import cycle protection**.

### How does GoArchTest compare to NetArchTest?

GoArchTest is inspired by NetArchTest but is adapted for Go's package-based architecture. While NetArchTest uses .NET's assembly and reflection capabilities, GoArchTest analyzes Go's AST (Abstract Syntax Tree) to understand code structure and dependencies.

### What Go versions are supported?

GoArchTest supports Go 1.18 and above.

### How does GoArchTest analyze Go code?

GoArchTest uses Go's modern type checking system and the `golang.org/x/tools/go/packages` package to analyze your code. It loads packages, extracts type information, and analyzes dependencies between types. Unlike older approaches that used the now-deprecated `ast.Package`, GoArchTest leverages Go's type system for more accurate and reliable analysis.

## Usage Questions

### How do I handle false positives in architecture tests?

If you have legitimate exceptions to your architecture rules, you can:

1. Make your rules more specific (e.g., use more precise namespace patterns)
2. Create separate tests for exceptional cases
3. Filter out specific false positives in your assertions

For example:

```go
// Get all types in the presentation layer
presentationTypes := types.That().ResideInNamespace("presentation").types

// Filter out specific exceptions
var nonExceptionTypes []*goarchtest.TypeInfo
for _, t := range presentationTypes {
    if t.Name != "LegacyException" {
        nonExceptionTypes = append(nonExceptionTypes, t)
    }
}

// Create a new TypeSet with the filtered types
filteredTypes := &goarchtest.TypeSet{
    types: nonExceptionTypes,
    originalTypes: presentationTypes,
}

// Now test the filtered types
result := filteredTypes.ShouldNot().HaveDependencyOn("data").GetResult()
```

### How can I test third-party dependencies?

You can use the `HaveDependencyOn` predicate to test for third-party dependencies:

```go
// Ensure data layer uses gorm
result := types.That().
    ResideInNamespace("data").
    Should().
    HaveDependencyOn("gorm.io").
    GetResult()

// Ensure we don't use deprecated libraries
result = types.That().
    ShouldNot().
    HaveDependencyOn("github.com/deprecated/package").
    GetResult()
```

### How do I handle generated code in architecture tests?

Generated code often violates architecture rules. You can exclude it from analysis by:

1. Placing generated code in specific directories and filtering them out
2. Using naming conventions for generated files and excluding them

```go
// Skip files in generated directories
projectPath, _ := filepath.Abs("./")
nonGeneratedPath := filepath.Join(projectPath, "src") // Only test non-generated code

types := goarchtest.InPath(nonGeneratedPath)
```

### Can I test microservices architecture?

Yes, you can test microservices by:

1. Testing each microservice independently
2. Testing the interactions between microservices

```go
// Test that services don't have direct dependencies on each other
result := types.That().
    ResideInNamespace("services/user").
    ShouldNot().
    HaveDependencyOn("services/product").
    GetResult()

// Test that services communicate through defined interfaces
result = types.That().
    ResideInNamespace("services").
    And().
    HaveDependencyOn("services").
    Should().
    HaveDependencyOn("proto").
    Or(types.That().ResideInNamespace("services").And().HaveDependencyOn("services").HaveDependencyOn("api/client")).
    GetResult()
```

## Technical Questions

### How does GoArchTest analyze code?

GoArchTest uses Go's built-in AST (Abstract Syntax Tree) parsing to analyze code structure. It examines:

1. Package declarations
2. Import statements
3. Type definitions (structs, interfaces)
4. Method and function signatures

### Does GoArchTest support generics?

Yes, GoArchTest can analyze code that uses Go generics (Go 1.18+).

### How can I debug architecture tests?

Add verbose logging to understand what types are being analyzed:

```go
// Print all types in a specific namespace
types := goarchtest.InPath(projectPath)
nsTypes := types.That().ResideInNamespace("mypackage").types

for _, t := range nsTypes {
    fmt.Printf("Type: %s in package %s\n", t.Name, t.Package)
    fmt.Printf("  Imports: %v\n", t.Imports)
}
```

### How can I contribute to GoArchTest?

Contributions are welcome! You can:

1. Submit bug reports and feature requests
2. Improve documentation
3. Submit pull requests with new features or bug fixes

See the [Contributing Guidelines](CONTRIBUTING.md) for more details.

## Performance Questions

### Is GoArchTest suitable for large codebases?

Yes, GoArchTest is designed to handle large codebases efficiently. However, for very large projects, you might want to:

1. Test specific directories rather than the entire codebase
2. Split architecture tests into multiple focused test functions
3. Run architecture tests separately from regular unit tests

### How can I improve the performance of architecture tests?

For large codebases, you can:

1. Be more specific in your namespace filters
2. Cache the Types instance between tests
3. Run architecture tests in parallel

```go
// Cache the Types instance
var cachedTypes *goarchtest.Types

func init() {
    projectPath, _ := filepath.Abs("./")
    cachedTypes = goarchtest.InPath(projectPath)
}

func TestArchitecture(t *testing.T) {
    // Use the cached types instance
    t.Run("DomainRules", func(t *testing.T) {
        t.Parallel()
        // Test domain rules
    })
    
    t.Run("ApplicationRules", func(t *testing.T) {
        t.Parallel()
        // Test application rules
    })
}
```
